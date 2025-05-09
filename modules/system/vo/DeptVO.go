package vo

type SysDeptDto struct {
	Id int64 `json:"id"`
	/** 节点名称 */
	Label string `json:"label"`
	/** 子节点 */
	Children []SysDeptDto `json:"children"`
}

type TreeSelect struct {
	Id       int64        `json:"id"`
	Label    string       `json:"label"`
	Children []TreeSelect `json:"children"`
}
