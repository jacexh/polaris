package protocols

type (
	// ReturnCode rest接口返回码
	ReturnCode int

	// AgentNode 基本的agent信息
	AgentNode struct {
		ID    string `json:"id"`
		Alias string `json:"alias,omitempty"`
	}

	// BaseRes 接口响应基本字段，适用于rest、ws接口
	BaseRes struct {
		Code    ReturnCode `json:"code"`
		Message string     `json:"message,omitempty"` // 当业务不成功时，该字段必定会返回
	}

	// WSMsgType Websocket 消息类型，用于区分不同的消息主体
	WSMsgType int

	// WSMsgDirection websocket消息流向，用于区分消息发起方
	WSMsgDirection int
)
