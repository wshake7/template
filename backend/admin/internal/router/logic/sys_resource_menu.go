package logic

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/slices_utils"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	MenuTypeCatalog  = "CATALOG"
	MenuTypeMenu     = "MENU"
	MenuTypeButton   = "BUTTON"
	MenuTypeEmbedded = "EMBEDDED"
	MenuTypeLink     = "LINK"
)

var validResourceMenuTypes = map[string]struct{}{
	MenuTypeCatalog:  {},
	MenuTypeMenu:     {},
	MenuTypeButton:   {},
	MenuTypeEmbedded: {},
	MenuTypeLink:     {},
}

type SysResourceMenuHandler struct{}

type RespSysResourceMenu struct {
	models.SysResourceMenu
	Children  []*RespSysResourceMenu `json:"children,omitempty"`
	CanWrite  bool                   `json:"canWrite"`
	CanDelete bool                   `json:"canDelete"`
}

type RespResourceMenuNode struct {
	ID          uint64                  `json:"id"`
	ParentID    *uint64                 `json:"parentID"`
	MenuType    string                  `json:"menuType"`
	Path        string                  `json:"path"`
	Redirect    string                  `json:"redirect"`
	Name        string                  `json:"name"`
	Component   string                  `json:"component"`
	Icon        string                  `json:"icon"`
	Order       int32                   `json:"order"`
	Hidden      bool                    `json:"hidden"`
	Authorities []string                `json:"authorities"`
	IsUrl       bool                    `json:"isUrl"`
	Children    []*RespResourceMenuNode `json:"children,omitempty"`
}

type ReqResourceMenuCreate struct {
	ParentID  *uint64           `json:"parentID"`
	MenuType  string            `json:"menuType" binding:"required" binding_msg:"required=菜单类型不能为空"`
	Path      string            `json:"path" binding:"max=1024" binding_msg:"max=路径最多1024位"`
	Redirect  string            `json:"redirect" binding:"max=1024" binding_msg:"max=重定向地址最多1024位"`
	Alias     string            `json:"alias" binding:"max=255" binding_msg:"max=路由别名最多255位"`
	Name      string            `json:"name" binding:"max=255" binding_msg:"max=路由命名最多255位"`
	Component string            `json:"component" binding:"max=255" binding_msg:"max=前端组件最多255位"`
	Metadata  datatypes.JSONMap `json:"metadata"`
	SortOrder int32             `json:"sortOrder"`
	IsEnabled bool              `json:"isEnabled"`
	Remark    string            `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceMenuUpdate struct {
	ID        uint64             `json:"id" binding:"required" binding_msg:"required=请求错误"`
	ParentID  *uint64            `json:"parentID"`
	MenuType  *string            `json:"menuType"`
	Path      *string            `json:"path" binding:"omitempty,max=1024" binding_msg:"max=路径最多1024位"`
	Redirect  *string            `json:"redirect" binding:"omitempty,max=1024" binding_msg:"max=重定向地址最多1024位"`
	Alias     *string            `json:"alias" binding:"omitempty,max=255" binding_msg:"max=路由别名最多255位"`
	Name      *string            `json:"name" binding:"omitempty,max=255" binding_msg:"max=路由命名最多255位"`
	Component *string            `json:"component" binding:"omitempty,max=255" binding_msg:"max=前端组件最多255位"`
	Metadata  *datatypes.JSONMap `json:"metadata"`
	SortOrder *int32             `json:"sortOrder"`
	IsEnabled *bool              `json:"isEnabled"`
	Remark    *string            `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceMenuBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择菜单资源,min=至少选择一项"`
}

// @Summary 获取菜单资源列表
// @Remark 分页查询菜单资源信息
// @Tags ResourceMenu
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespSysResourceMenu]} "成功"
// @Router /api/sys/resource/menu/list [post]
func (*SysResourceMenuHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespSysResourceMenu], error) {
	if req.GetOrderBy() == "" {
		orderBy := "sort_order asc,id asc"
		req.OrderBy = &orderBy
	}
	pagination, err := query.SysResourceMenu.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}

	items := buildResourceMenuRespTree(pagination.Items)
	return &gormc.PagingResult[RespSysResourceMenu]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

// @Summary 获取当前用户动态菜单树
// @Remark 查询启用的目录、菜单、内嵌页和外链菜单
// @Tags ResourceMenu
// @Produce json
// @Success 200 {object} res.Response{data=[]RespResourceMenuNode} "成功"
// @Router /api/sys/resource/menu/tree [get]
func (*SysResourceMenuHandler) Tree(ctx *handler.Ctx) (*[]*RespResourceMenuNode, error) {
	sysResourceMenu := query.SysResourceMenu
	items, err := sysResourceMenu.
		Where(
			sysResourceMenu.IsEnabled.Is(true),
			sysResourceMenu.MenuType.In(MenuTypeCatalog, MenuTypeMenu, MenuTypeEmbedded, MenuTypeLink),
		).
		Order(sysResourceMenu.SortOrder.Asc(), sysResourceMenu.ID.Asc()).
		Find()
	if err != nil {
		return nil, res.FailDefault
	}
	tree := buildResourceMenuNodeTree(items)
	return &tree, nil
}

// @Summary 创建菜单资源
// @Remark 创建新的菜单资源节点
// @Tags ResourceMenu
// @Accept json
// @Produce json
// @Param req body ReqResourceMenuCreate true "菜单资源创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/menu/create [post]
func (*SysResourceMenuHandler) Create(ctx *handler.Ctx, req *ReqResourceMenuCreate) error {
	req.normalize()
	if err := validateResourceMenuValues(req.MenuType, req.Name, req.Path, req.Component); err != nil {
		return err
	}

	operationID := ctx.SessionInfo.Id
	return query.Q.Transaction(func(tx *query.Query) error {
		parentID, err := normalizeResourceMenuParentID(tx, req.MenuType, req.ParentID)
		if err != nil {
			return err
		}

		item := &models.SysResourceMenu{
			OperatorID: mixin.OperatorID{
				CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
				UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
			},
			Remark:    mixin.Remark{Remark: req.Remark},
			SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
			Metadata:  mixin.Metadata{Metadata: req.Metadata},
			IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
			MenuType:  req.MenuType,
			Path:      req.Path,
			Redirect:  req.Redirect,
			Alias:     req.Alias,
			Name:      req.Name,
			Component: req.Component,
			ParentID:  parentID,
		}
		if err := tx.SysResourceMenu.Create(item); err != nil {
			return res.FailDefault
		}
		return updateResourceMenuTreePath(tx, item.ID)
	})
}

// @Summary 更新菜单资源
// @Remark 根据 ID 更新菜单资源；启停状态也通过该接口更新
// @Tags ResourceMenu
// @Accept json
// @Produce json
// @Param req body ReqResourceMenuUpdate true "菜单资源更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/menu/update [post]
func (*SysResourceMenuHandler) Update(ctx *handler.Ctx, req *ReqResourceMenuUpdate) error {
	req.normalize()
	return query.Q.Transaction(func(tx *query.Query) error {
		sysResourceMenu := tx.SysResourceMenu
		current, err := sysResourceMenu.Where(sysResourceMenu.ID.Eq(req.ID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("菜单资源不存在")
			}
			return res.FailDefault
		}

		menuType := current.MenuType
		if req.MenuType != nil {
			menuType = strings.TrimSpace(*req.MenuType)
		}
		name := current.Name
		if req.Name != nil {
			name = strings.TrimSpace(*req.Name)
		}
		path := current.Path
		if req.Path != nil {
			path = strings.TrimSpace(*req.Path)
		}
		component := current.Component
		if req.Component != nil {
			component = strings.TrimSpace(*req.Component)
		}
		if err := validateResourceMenuValues(menuType, name, path, component); err != nil {
			return err
		}

		updates := map[string]any{
			"updated_by": ctx.SessionInfo.Id,
		}
		parentID := current.ParentID
		parentChanged := false
		if req.ParentID != nil {
			var err error
			parentID, err = normalizeResourceMenuParentID(tx, menuType, req.ParentID)
			if err != nil {
				return err
			}
			if parentID != nil && *parentID == req.ID {
				return res.FailMsg("父级不能选择自身")
			}
			if err := ensureResourceMenuParentNotDescendant(tx, req.ID, parentID); err != nil {
				return err
			}
			updates["parent_id"] = parentID
			parentChanged = !sameOptionalUint64(current.ParentID, parentID)
		}
		if err := validateResourceMenuParentType(tx, menuType, parentID); err != nil {
			return err
		}
		if req.MenuType != nil {
			updates["menu_type"] = menuType
		}
		if req.Path != nil {
			updates["path"] = path
		}
		if req.Redirect != nil {
			updates["redirect"] = strings.TrimSpace(*req.Redirect)
		}
		if req.Alias != nil {
			updates["alias"] = strings.TrimSpace(*req.Alias)
		}
		if req.Name != nil {
			updates["name"] = name
		}
		if req.Component != nil {
			updates["component"] = strings.TrimSpace(*req.Component)
		}
		if req.Metadata != nil {
			updates["metadata"] = *req.Metadata
		}
		if req.SortOrder != nil {
			updates["sort_order"] = *req.SortOrder
		}
		if req.IsEnabled != nil {
			updates["is_enabled"] = *req.IsEnabled
		}
		if req.Remark != nil {
			updates["remark"] = strings.TrimSpace(*req.Remark)
		}

		info, err := sysResourceMenu.Where(sysResourceMenu.ID.Eq(req.ID)).Updates(updates)
		if err != nil {
			return res.FailDefault
		}
		if info.RowsAffected == 0 {
			return res.FailMsg("菜单资源不存在")
		}
		if parentChanged {
			return updateResourceMenuTreePath(tx, req.ID)
		}
		return nil
	})
}

// @Summary 删除菜单资源
// @Remark 根据 ID 列表批量删除菜单资源；存在子节点时拒绝删除
// @Tags ResourceMenu
// @Accept json
// @Produce json
// @Param req body ReqResourceMenuBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/menu/del [post]
func (*SysResourceMenuHandler) Del(ctx *handler.Ctx, req *ReqResourceMenuBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	if len(ids) == 0 {
		return res.FailMsg("请选择菜单资源")
	}
	sysResourceMenu := query.SysResourceMenu
	children, err := sysResourceMenu.
		Select(sysResourceMenu.ID).
		Where(sysResourceMenu.ParentID.In(ids...)).
		Find()
	if err != nil {
		return res.FailDefault
	}
	if len(children) > 0 {
		return res.FailMsg("存在子节点，不能删除")
	}

	info, err := sysResourceMenu.Where(sysResourceMenu.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("菜单资源不存在")
	}
	return nil
}

func (req *ReqResourceMenuCreate) normalize() {
	req.MenuType = strings.ToUpper(strings.TrimSpace(req.MenuType))
	req.Path = strings.TrimSpace(req.Path)
	req.Redirect = strings.TrimSpace(req.Redirect)
	req.Alias = strings.TrimSpace(req.Alias)
	req.Name = strings.TrimSpace(req.Name)
	req.Component = strings.TrimSpace(req.Component)
	req.Remark = strings.TrimSpace(req.Remark)
	if req.Metadata == nil {
		req.Metadata = datatypes.JSONMap{}
	}
	normalizeResourceMenuMetadata(req.Metadata)
}

func (req *ReqResourceMenuUpdate) normalize() {
	trimStringPtr(req.MenuType, func(v string) string { return strings.ToUpper(v) })
	trimStringPtr(req.Path, nil)
	trimStringPtr(req.Redirect, nil)
	trimStringPtr(req.Alias, nil)
	trimStringPtr(req.Name, nil)
	trimStringPtr(req.Component, nil)
	trimStringPtr(req.Remark, nil)
	if req.Metadata != nil {
		if *req.Metadata == nil {
			*req.Metadata = datatypes.JSONMap{}
		}
		normalizeResourceMenuMetadata(*req.Metadata)
	}
}

func trimStringPtr(value *string, transform func(string) string) {
	if value == nil {
		return
	}
	next := strings.TrimSpace(*value)
	if transform != nil {
		next = transform(next)
	}
	*value = next
}

func normalizeResourceMenuMetadata(metadata datatypes.JSONMap) {
	if metadata == nil {
		return
	}
	if order, ok := metadata["order"]; ok {
		metadata["order"] = int32FromAny(order)
	}
	if hidden, ok := metadata["hidden"]; ok {
		metadata["hidden"] = boolFromAny(hidden)
	}
	if authorities, ok := metadata["authorities"]; ok {
		metadata["authorities"] = stringsFromAny(authorities)
	}
}

func validateResourceMenuValues(menuType, name, path, component string) error {
	if _, ok := validResourceMenuTypes[menuType]; !ok {
		return res.FailMsg("菜单类型无效")
	}
	if strings.TrimSpace(name) == "" {
		return res.FailMsg("菜单名称不能为空")
	}
	if (menuType == MenuTypeMenu || menuType == MenuTypeEmbedded || menuType == MenuTypeLink) && strings.TrimSpace(path) == "" {
		return res.FailMsg("菜单路径不能为空")
	}
	if menuType == MenuTypeMenu && strings.TrimSpace(component) == "" {
		return res.FailMsg("组件路径不能为空")
	}
	return nil
}

func normalizeResourceMenuParentID(tx *query.Query, menuType string, parentID *uint64) (*uint64, error) {
	if parentID == nil || *parentID == 0 {
		if menuType == MenuTypeCatalog {
			return nil, nil
		}
		return nil, res.FailMsg("上级菜单不能为空")
	}
	if err := validateResourceMenuParentType(tx, menuType, parentID); err != nil {
		return nil, err
	}
	return parentID, nil
}

func validateResourceMenuParentType(tx *query.Query, menuType string, parentID *uint64) error {
	if parentID == nil || *parentID == 0 {
		if menuType == MenuTypeCatalog {
			return nil
		}
		return res.FailMsg("上级菜单不能为空")
	}
	sysResourceMenu := tx.SysResourceMenu
	parent, err := sysResourceMenu.
		Select(sysResourceMenu.ID, sysResourceMenu.MenuType).
		Where(sysResourceMenu.ID.Eq(*parentID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("父级菜单不存在")
		}
		return res.FailDefault
	}
	if (menuType == MenuTypeCatalog || menuType == MenuTypeMenu) && parent.MenuType != MenuTypeCatalog {
		return res.FailMsg("目录和菜单的上级菜单只能是目录")
	}
	if (menuType == MenuTypeButton || menuType == MenuTypeEmbedded || menuType == MenuTypeLink) && parent.MenuType != MenuTypeMenu {
		return res.FailMsg("按钮、内嵌和外链的上级菜单只能是菜单")
	}
	return nil
}

func ensureResourceMenuParentNotDescendant(tx *query.Query, id uint64, parentID *uint64) error {
	if parentID == nil {
		return nil
	}
	parent, err := tx.SysResourceMenu.
		Select(tx.SysResourceMenu.ID, tx.SysResourceMenu.TreePath).
		Where(tx.SysResourceMenu.ID.Eq(*parentID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("父级菜单不存在")
		}
		return res.FailDefault
	}
	if strings.Contains(parentTreePath(parent), fmt.Sprintf("/%d/", id)) {
		return res.FailMsg("父级不能选择自身或子节点")
	}
	return nil
}

func updateResourceMenuTreePath(tx *query.Query, id uint64) error {
	sysResourceMenu := tx.SysResourceMenu
	current, err := sysResourceMenu.
		Select(sysResourceMenu.ID, sysResourceMenu.ParentID, sysResourceMenu.TreePath).
		Where(sysResourceMenu.ID.Eq(id)).
		First()
	if err != nil {
		return err
	}

	treePath := fmt.Sprintf("/%d/", current.ID)
	if current.ParentID != nil {
		parent, err := sysResourceMenu.
			Select(sysResourceMenu.ID, sysResourceMenu.TreePath).
			Where(sysResourceMenu.ID.Eq(*current.ParentID)).
			First()
		if err != nil {
			return err
		}
		treePath = parentTreePath(parent) + strconv.FormatUint(current.ID, 10) + "/"
	}

	_, err = sysResourceMenu.
		Where(sysResourceMenu.ID.Eq(current.ID)).
		Update(sysResourceMenu.TreePath, treePath)
	if err != nil {
		return err
	}

	children, err := sysResourceMenu.
		Select(sysResourceMenu.ID).
		Where(sysResourceMenu.ParentID.Eq(current.ID)).
		Find()
	if err != nil {
		return err
	}
	for _, child := range children {
		if err := updateResourceMenuTreePath(tx, child.ID); err != nil {
			return err
		}
	}
	return nil
}

func parentTreePath(item *models.SysResourceMenu) string {
	if item.TreePath != nil && *item.TreePath != "" {
		return *item.TreePath
	}
	return fmt.Sprintf("/%d/", item.ID)
}

func buildResourceMenuRespTree(items []*models.SysResourceMenu) []*RespSysResourceMenu {
	nodes := make(map[uint64]*RespSysResourceMenu, len(items))
	roots := make([]*RespSysResourceMenu, 0)
	for _, item := range items {
		if item == nil {
			continue
		}
		nodes[item.ID] = &RespSysResourceMenu{
			SysResourceMenu: *item,
			CanWrite:        true,
			CanDelete:       true,
		}
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		node := nodes[item.ID]
		if item.ParentID != nil {
			if parent, ok := nodes[*item.ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	sortResourceMenuRespNodes(roots)
	return roots
}

func sortResourceMenuRespNodes(nodes []*RespSysResourceMenu) {
	slices.SortFunc(nodes, func(a, b *RespSysResourceMenu) int {
		if a.SortOrder.SortOrder != b.SortOrder.SortOrder {
			return int(a.SortOrder.SortOrder - b.SortOrder.SortOrder)
		}
		return int(a.ID - b.ID)
	})
	for _, node := range nodes {
		sortResourceMenuRespNodes(node.Children)
	}
}

func buildResourceMenuNodeTree(items []*models.SysResourceMenu) []*RespResourceMenuNode {
	nodes := make(map[uint64]*RespResourceMenuNode, len(items))
	roots := make([]*RespResourceMenuNode, 0)
	for _, item := range items {
		if item == nil {
			continue
		}
		node := resourceMenuToNode(item)
		if node.Hidden {
			continue
		}
		nodes[item.ID] = node
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		node, ok := nodes[item.ID]
		if !ok {
			continue
		}
		if item.ParentID != nil {
			if parent, ok := nodes[*item.ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	sortResourceMenuNodeNodes(roots)
	return roots
}

func resourceMenuToNode(item *models.SysResourceMenu) *RespResourceMenuNode {
	metadata := item.Metadata.Metadata
	return &RespResourceMenuNode{
		ID:          item.ID,
		ParentID:    item.ParentID,
		MenuType:    item.MenuType,
		Path:        item.Path,
		Redirect:    item.Redirect,
		Name:        item.Name,
		Component:   item.Component,
		Icon:        stringFromAny(metadata["icon"]),
		Order:       int32FromAny(metadata["order"]),
		Hidden:      boolFromAny(metadata["hidden"]),
		Authorities: stringsFromAny(metadata["authorities"]),
		IsUrl:       item.MenuType == MenuTypeLink,
	}
}

func sortResourceMenuNodeNodes(nodes []*RespResourceMenuNode) {
	slices.SortFunc(nodes, func(a, b *RespResourceMenuNode) int {
		if a.Order != b.Order {
			return int(a.Order - b.Order)
		}
		return int(a.ID - b.ID)
	})
	for _, node := range nodes {
		sortResourceMenuNodeNodes(node.Children)
	}
}

func sameOptionalUint64(a, b *uint64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func stringsFromAny(value any) []string {
	switch v := value.(type) {
	case []string:
		return compactStrings(v)
	case []any:
		items := make([]string, 0, len(v))
		for _, item := range v {
			if text := stringFromAny(item); text != "" {
				items = append(items, text)
			}
		}
		return compactStrings(items)
	case string:
		return compactStrings(strings.Split(v, ","))
	default:
		return []string{}
	}
}

func compactStrings(values []string) []string {
	items := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}
	return items
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	case float64:
		return v != 0
	case int:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	default:
		return false
	}
}

func int32FromAny(value any) int32 {
	switch v := value.(type) {
	case int32:
		return v
	case int:
		return int32(v)
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		parsed, _ := strconv.ParseInt(v, 10, 32)
		return int32(parsed)
	default:
		return 0
	}
}
