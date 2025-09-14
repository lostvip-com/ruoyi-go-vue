package vo_emqx

// EmqxOnlineVO emqx固定格式，没有经过自定义
type EmqxOnlineVO struct {
	Node      string `json:"node"`
	Reason    string `json:"reason"`
	Timestamp int    `json:"timestamp"`
	Peername  string `json:"peername"`
	Sockname  string `json:"sockname"`
	Metadata  struct {
		RuleId string `json:"rule_id"`
	} `json:"metadata"`
	Event     string `json:"event"`
	Username  string `json:"username"`
	EventType string `json:"event_type"`
	ProtoVer  int    `json:"proto_ver"`
	Clientid  string `json:"clientid"`
}
