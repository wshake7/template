package logic

import "C"
import (
	"admin/auth"
	"admin/fiberc/handler"
	"admin/services/orm/query"
	"errors"
	"github.com/click33/sa-token-go/stputil"
	"go-common/utils/encrypt/rsa_util"
	"go-common/utils/passwd"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccountHandler struct{}

type ReqAccountPwdLogin struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
}

type ResAccountPwdLogin struct {
	Token     string `json:"token"`
	PublicKey string `json:"publicKey"`
}

func (*AccountHandler) PwdLogin(ctx *handler.Ctx, req *ReqAccountPwdLogin) (*ResAccountPwdLogin, error) {
	sysUser := query.SysUser
	var result struct {
		ID       uint64
		Password string
	}
	err := sysUser.Where(sysUser.Username.Eq(req.Username)).Select(sysUser.ID, sysUser.Password).Scan(&result)
	logger := ctx.L().With(zap.String("username", req.Username))
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
	Token string `cookie:"token"`
}

func (*AccountHandler) Logout(ctx *handler.Ctx, req *ReqAccountLogout) error {
	loginID, err := stputil.GetLoginID(req.Token)
	if err != nil {
		return errors.New("操作失败")
	}
	err = stputil.Logout(loginID)
	if err != nil {
		ctx.L().Error("退出登录失败", zap.Error(err), zap.String("token", req.Token))
		return errors.New("操作失败")
	}
	return nil
}

type ReqAccountChangePwd struct {
	OldPwd string `json:"oldPwd"`
	NewPwd string `json:"newPwd"`
}

type ChangeTest struct {
	Id       uint64
	Nickname string `change:"昵称"`
}

func (*ChangeTest) ChangeString(before *ChangeTest) string {
	return "你好"
}

func (*AccountHandler) ChangePwdQuery(ctx *handler.Ctx, test *ChangeTest) (*ChangeTest, error) {
	var c ChangeTest
	sysUser := query.SysUser
	err := sysUser.Where(sysUser.ID.Eq(1)).Scan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (*AccountHandler) ChangePwd(ctx *handler.Ctx, req *ReqAccountChangePwd) error {
	sysUser := query.SysUser
	_, err := sysUser.Where(sysUser.ID.Eq(1)).Update(sysUser.Nickname, "11234")
	return err
}
