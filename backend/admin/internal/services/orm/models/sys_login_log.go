package models

import (
	"time"

	"orm-crud/gormc/mixin"

	"gorm.io/gorm"
)

func init() {
	Models = append(Models, &SysLoginLog{})
}

type SysLoginLog struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	LoginIP        string     `gorm:"column:login_ip;type:varchar(64);default:'';comment:登录IP地址"`
	LoginMAC       string     `gorm:"column:login_mac;type:varchar(128);default:'';comment:登录MAC地址"`
	LoginTime      *time.Time `gorm:"column:login_time;comment:登录时间"`
	UserAgent      string     `gorm:"column:user_agent;type:text;default:'';comment:浏览器的用户代理信息"`
	BrowserName    string     `gorm:"column:browser_name;type:varchar(128);default:'';comment:浏览器名称"`
	BrowserVersion string     `gorm:"column:browser_version;type:varchar;default:''(128);default:'';comment:浏览器版本"`
	ClientID       string     `gorm:"column:client_id;type:varchar(128);comment:客户端ID"`
	ClientName     string     `gorm:"column:client_name;type:varchar(128);default:'';comment:客户端名称"`
	OSName         string     `gorm:"column:os_name;type:varchar(128);default:'';comment:操作系统名称"`
	OSVersion      string     `gorm:"column:os_version;type:varchar(128);default:'';comment:操作系统版本"`
	SysUserID      uint64     `gorm:"column:sys_user_id;type:bigint;comment:操作者用户ID" json:"sysUserID"`
	StatusCode     int32      `gorm:"column:status_code;default:0;comment:状态码"`
	Success        bool       `gorm:"column:success;default:false;comment:操作成功"`
	Reason         string     `gorm:"column:reason;type:varchar(255);default:'';comment:登录失败原因"`
	Location       string     `gorm:"column:location;type:varchar(255);default:'';comment:登录地理位置"`

	SysUser *SysUser `gorm:"foreignKey:SysUserID;references:ID" json:"sysUser"`
}

func (m *SysLoginLog) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedAt.CreatedAt == nil {
		now := time.Now()
		m.CreatedAt.CreatedAt = &now
	}
	return nil
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
