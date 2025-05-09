package vo

type UserRoleReq struct {
	RoleId string `json:"roleId" form:"roleId"`
	UserId string `json:"userId" form:"userId"`
}
