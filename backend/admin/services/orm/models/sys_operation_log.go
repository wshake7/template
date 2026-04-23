package models

import (
	"orm-crud/gorm/mixin"
)

func init() {
	Models = append(Models, &SysOperationLog{})
}

type SysOperationLog struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	RequestID      string `gorm:"column:request_id;type:varchar(128);not null;comment:请求ID"`
	Method         string `gorm:"column:method;type:varchar(16);not null;comment:请求方法"`
	Module         string `gorm:"column:module;type:varchar(255);default:'';comment:模块"`
	Path           string `gorm:"column:path;type:varchar(255);not null;comment:请求路径"`
	Referer        string `gorm:"column:referer;type:varchar(255);default:'';comment:请求源"`
	BeforeChange   string `gorm:"column:before_change;type:text;default:'';comment:变更前内容"`
	AfterChange    string `gorm:"column:after_change;type:text;default:'';comment:变更后内容"`
	FormatChange   string `gorm:"column:format_change;type:text;default:'';comment:格式化变化内容"`
	RequestURI     string `gorm:"column:request_uri;type:varchar(255);default:'';comment:请求URI"`
	RequestBody    string `gorm:"column:request_body;type:text;default:'';comment:请求体"`
	RequestHeader  string `gorm:"column:request_header;type:text;default:'';comment:请求头"`
	Response       string `gorm:"column:response;type:text;default:'';comment:响应信息"`
	CostTime       int64  `gorm:"column:cost_time;type:bigint;default:0;comment:操作耗时"`
	UserID         uint64 `gorm:"column:user_id;type:bigint;default:0;comment:操作者用户ID"`
	Username       string `gorm:"column:username;type:varchar(128);default:'';comment:操作者账号名"`
	ClientIP       string `gorm:"column:client_ip;type:varchar(64);default:'';comment:操作者IP"`
	StatusCode     int    `gorm:"column:status_code;type:int;default:0;comment:状态码"`
	Reason         string `gorm:"column:reason;type:varchar(255);default:'';comment:操作失败原因"`
	Success        bool   `gorm:"column:success;type:boolean;default:false;comment:操作成功"`
	Location       string `gorm:"column:location;type:varchar(255);default:'';comment:操作地理位置"`
	UserAgent      string `gorm:"column:user_agent;type:text;default:'';comment:浏览器的用户代理信息"`
	BrowserName    string `gorm:"column:browser_name;type:varchar(128);default:'';comment:浏览器名称"`
	BrowserVersion string `gorm:"column:browser_version;type:varchar(128);default:'';comment:浏览器版本"`
	ClientID       string `gorm:"column:client_id;type:varchar(128);default:'';comment:客户端ID"`
	ClientName     string `gorm:"column:client_name;type:varchar(128);default:'';comment:客户端名称"`
	OSName         string `gorm:"column:os_name;type:varchar(128);default:'';comment:操作系统名称"`
	OSVersion      string `gorm:"column:os_version;type:varchar(128);default:'';comment:操作系统版本"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
