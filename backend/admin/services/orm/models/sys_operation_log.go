package models

import (
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysOperationLog{})
}

type SysOperationLog struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	RequestID      string `gorm:"column:request_id;type:varchar(128);not null;comment:请求ID" json:"requestID,omitempty"`
	Method         string `gorm:"column:method;type:varchar(16);not null;comment:请求方法" json:"method,omitempty"`
	Module         string `gorm:"column:module;type:varchar(255);default:'';comment:模块" json:"module,omitempty"`
	Path           string `gorm:"column:path;type:varchar(255);not null;comment:请求路径" json:"path,omitempty"`
	Referer        string `gorm:"column:referer;type:varchar(255);default:'';comment:请求源" json:"referer,omitempty"`
	BeforeChange   string `gorm:"column:before_change;type:text;default:'';comment:变更前内容" json:"beforeChange,omitempty"`
	AfterChange    string `gorm:"column:after_change;type:text;default:'';comment:变更后内容" json:"afterChange,omitempty"`
	FormatChange   string `gorm:"column:format_change;type:text;default:'';comment:格式化变化内容" json:"formatChange,omitempty"`
	RequestURI     string `gorm:"column:request_uri;type:varchar(255);default:'';comment:请求URI" json:"requestURI,omitempty"`
	RequestBody    string `gorm:"column:request_body;type:text;default:'';comment:请求体" json:"requestBody,omitempty"`
	RequestHeader  string `gorm:"column:request_header;type:text;default:'';comment:请求头" json:"requestHeader,omitempty"`
	Response       string `gorm:"column:response;type:text;default:'';comment:响应信息" json:"response,omitempty"`
	CostTime       int64  `gorm:"column:cost_time;type:bigint;default:0;comment:操作耗时" json:"costTime,omitempty"`
	UserID         uint64 `gorm:"column:user_id;type:bigint;default:0;comment:操作者用户ID" json:"userID,omitempty"`
	Username       string `gorm:"column:username;type:varchar(128);default:'';comment:操作者账号名" json:"username,omitempty"`
	ClientIP       string `gorm:"column:client_ip;type:varchar(64);default:'';comment:操作者IP" json:"clientIP,omitempty"`
	StatusCode     int    `gorm:"column:status_code;type:int;default:0;comment:状态码" json:"statusCode,omitempty"`
	Reason         string `gorm:"column:reason;type:varchar(255);default:'';comment:操作失败原因" json:"reason,omitempty"`
	Success        bool   `gorm:"column:success;type:boolean;default:false;comment:操作成功" json:"success,omitempty"`
	Location       string `gorm:"column:location;type:varchar(255);default:'';comment:操作地理位置" json:"location,omitempty"`
	UserAgent      string `gorm:"column:user_agent;type:text;default:'';comment:浏览器的用户代理信息" json:"userAgent,omitempty"`
	BrowserName    string `gorm:"column:browser_name;type:varchar(128);default:'';comment:浏览器名称" json:"browserName,omitempty"`
	BrowserVersion string `gorm:"column:browser_version;type:varchar(128);default:'';comment:浏览器版本" json:"browserVersion,omitempty"`
	ClientID       string `gorm:"column:client_id;type:varchar(128);default:'';comment:客户端ID" json:"clientID,omitempty"`
	ClientName     string `gorm:"column:client_name;type:varchar(128);default:'';comment:客户端名称" json:"clientName,omitempty"`
	OSName         string `gorm:"column:os_name;type:varchar(128);default:'';comment:操作系统名称" json:"oSName,omitempty"`
	OSVersion      string `gorm:"column:os_version;type:varchar(128);default:'';comment:操作系统版本" json:"oSVersion,omitempty"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
