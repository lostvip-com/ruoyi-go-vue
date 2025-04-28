package vo

// RouterVO 菜单权限
type RouterVO struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Hidden     bool       `json:"hidden"`
	Redirect   string     `json:"redirect"`
	Component  string     `json:"component"`
	AlwaysShow bool       `json:"alwaysShow"`
	Children   []RouterVO `json:"children"`
	Meta       Meta       `json:"meta"`

	MenuId   int64 `json:"menuId"` // 菜单ID
	ParentId int64 `json:"parentId"`
}

type Meta struct {
	Title   string `json:"title"`
	Icon    string `json:"icon"`
	NoCache bool   `json:"noCache"`
	Link    string `json:"link"`
}
