package vo_emqx

// emqx 中访问控制->客户端认证->HTTP认证
type EmqxLoginVO struct {
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
}
