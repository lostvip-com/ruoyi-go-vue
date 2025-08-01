package vo

type UserRoleReq struct {
	RoleId string `json:"roleId" form:"roleId"`
	UserId string `json:"userId" form:"userId"`
}

type RoleStatusReq struct {
	RoleId int    `json:"roleId" form:"roleId"`
	Status string `json:"status" form:"status"`
}
