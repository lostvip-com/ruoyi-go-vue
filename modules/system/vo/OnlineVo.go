package vo

type OnlineVo struct {
	SessionId      string `json:"sessionId"`
	UserName       string `json:"UserName"`
	DeptName       string `json:"deptName"`
	Ipaddr         string `json:"ipaddr"`
	LoginLocation  string `json:"loginLocation"`
	Browser        string `json:"browser"`
	Os             string `json:"os"`
	Status         string `json:"status"`
	StartTimestamp string `json:"startTimestamp"`
	LastAccessTime string `json:"lastAccessTime"`
	CreateTime     string `json:"createTime"`
}

// 增
func (e *OnlineVo) Save() error {
	return nil
}

// 查
func (e *OnlineVo) FindById() error {
	return nil
}

// 改
func (e *OnlineVo) Update() error {
	return nil
}

// 删
func (e *OnlineVo) Delete() error {
	return nil
}
