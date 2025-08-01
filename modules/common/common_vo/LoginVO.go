package common_vo

// LoginInfo 用户登陆信息
type LoginInfo struct {
	UserId   int    `json:"userId" `
	DeptId   int    `json:"deptId"`   //结束时间
	Username string `json:"username"` //登录名
	RoleKeys string `json:"roleKeys"` //登录名
	Avatar   string `json:"avatar"`   //登录名

	TenantId int `json:"tenantId"`
}
