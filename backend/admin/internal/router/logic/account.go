package logic

import "C"
import (
	"admin/internal/auth"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm"
	"admin/internal/services/orm/query"
	"admin/internal/services/orm/repo"
	"errors"
	"go-common/utils/encrypt/rsa_util"
	"go-common/utils/passwd"

	"github.com/click33/sa-token-go/stputil"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccountHandler struct{}

type ReqAccountPwdLogin struct {
	Username string `json:"username" binding:"required,max=24" binding_msg:"required=用户名不能为空,max=用户名最多24位"`
	Pwd      string `json:"pwd" binding:"required,min=6" binding_msg:"required=密码不能为空,min=密码最少6位"`
}

type ResAccountPwdLogin struct {
	Token     string `json:"token"`
	PublicKey string `json:"publicKey"`
}

// @Summary 用户密码登录
// @Description 通过用户名和密码登录系统，并返回 token 与会话公钥
// @Tags Account
// @Accept json
// @Produce json
// @Param req body ReqAccountPwdLogin true "登录请求"
// @Success 200 {object} res.Response{data=ResAccountPwdLogin} "成功"
// @Router /api/account/login/pwd [post]
func (*AccountHandler) PwdLogin(ctx *handler.Ctx, req *ReqAccountPwdLogin) (*ResAccountPwdLogin, error) {
	logger := ctx.L().With(zap.String("username", req.Username))
	sysUser := query.SysUser
	result, err := repo.SysUserRepo.Get(ctx.Context(), orm.DB(), sysUser.ID, sysUser.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码无效")
		}
		logger.Error("获取用户失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 校验密码
	if !passwd.Match(req.Pwd, result.Password) {
		return nil, errors.New("用户名或密码无效")
	}
	token, err := stputil.Login(result.ID)
	if err != nil {
		logger.Error("获取token失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	session, err := auth.GetSession(result.ID)
	if err != nil {
		logger.Error("获取session失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	privateKey, publicKey, err := rsa_util.GenerateKeyPair()
	if err != nil {
		logger.Error("获取rsaKey错误", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	err = session.SaveInfo(&auth.SessionInfo{
		PrivateKey: privateKey,
		Id:         result.ID,
	})
	if err != nil {
		logger.Error("保存SessionInfo错误", zap.Error(err))
		return nil, errors.New("登录失败")
	}
	return &ResAccountPwdLogin{
		Token:     token,
		PublicKey: publicKey,
	}, nil
}

type ReqAccountLogout struct {
	Token string `cookie:"token" binding:"required" binding_msg:"required=请求错误'"`
}

// @Summary 退出登录
// @Description 使当前登录态失效
// @Tags Account
// @Produce json
// @Param token header string true "登录 token"
// @Success 200 {object} res.Response "成功"
// @Router /api/account/logout [get]
func (*AccountHandler) Logout(ctx *handler.Ctx, req *ReqAccountLogout) error {
	loginID, err := stputil.GetLoginID(req.Token)
	if err != nil {
		ctx.L().Error("获取loginId失败", zap.Error(err))
		return auth.CheckLoginErr(err)
	}
	err = stputil.Logout(loginID)
	if err != nil {
		ctx.L().Error("退出登录失败", zap.Error(err), zap.String("token", req.Token))
		return auth.CheckLoginErr(err)
	}
	return nil
}

type ReqAccountChangePwd struct {
	OldPwd string `json:"oldPwd" binding:"required,min=6" binding_msg:"required=原始密码不能为空,min=原始密码最少6位"`
	NewPwd string `json:"newPwd" binding:"required,min=6" binding_msg:"required=新密码不能为空,min=新密码最少6位"`
}

// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags Account
// @Accept json
// @Produce json
// @Param req body ReqAccountChangePwd true "修改密码请求"
// @Success 200 {object} res.Response "成功"
// @Router /api/account/changePwd [post]
func (*AccountHandler) ChangePwd(ctx *handler.Ctx, req *ReqAccountChangePwd) error {
	info := ctx.SessionInfo
	sysUser := query.SysUser
	result, err := repo.SysUserRepo.Get(ctx.Context(), orm.DB().Where(sysUser.ID.Eq(info.Id)), sysUser.Password)
	if err != nil {
		ctx.L().Error("获取用户密码失败", zap.Error(err))
		return res.FailDefault
	}

	if !passwd.Match(req.NewPwd, result.Password) {
		return errors.New("原密码错误")
	}

	encodePwd, err := passwd.Encode(req.NewPwd)
	if err != nil {
		ctx.L().Error("密码加密失败", zap.Error(err))
		return res.FailDefault
	}
	sysUser = query.SysUser
	_, err = repo.SysUserRepo.UpdateMap(map[string]any{
		sysUser.Password.ColumnName().String(): encodePwd,
	}, sysUser.ID.Eq(info.Id))
	if err != nil {
		ctx.L().Error("修改密码失败", zap.Error(err))
		return res.FailDefault
	}

	return err
}
