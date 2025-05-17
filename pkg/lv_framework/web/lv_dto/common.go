package lv_dto

type BunissType int

// 业务类型
const (
	Buniss_Other BunissType = 0 //0其它
	Buniss_Add   BunissType = 1 //1新增
	Buniss_Edit  BunissType = 2 //2修改
	Buniss_Del   BunissType = 3 //3删除
)

// 响应结果
const (
	SUCCESS      = 200 // 成功
	ERROR        = 500 //错误
	UNAUTHORIZED = 403 //无权限
	FAIL         = -1  //失败
)

// 通用api响应
type CommonRes struct {
	Code int         `json:"code"` //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg  string      `json:"msg"`  //消息
	Data interface{} `json:"data"` //数据内容
}

// 验证码响应
type CaptchaRes struct {
	Code           int         `json:"code"` //响应编码 0 成功 500 错误 403 无权限
	Msg            string      `json:"msg"`  //消息
	Img            interface{} `json:"img"`  //数据内容
	Uuid           string      `json:"uuid"` //验证码ID
	CaptchaEnabled bool        `json:"captchaEnabled"`
}

// 通用分页表格响应
type TableDataInfo struct {
	Total any         `json:"total"` //总数
	Rows  interface{} `json:"rows"`  //数据
	Code  int         `json:"code"`  //响应编码 200 成功 500 错误 403 无权限
	Msg   string      `json:"msg"`   //消息
}
