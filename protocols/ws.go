package protocols

import "encoding/json"

const (
	FromClient WSMsgDirection = iota
	FromServer

	WSMTConnection    WSMsgType = 10000 // ws连接类消息
	WSMTDisconnection WSMsgType = 10002 // ws主动断链消息

	WSMTTask             WSMsgType = 20000 // ws分发任务消息，console下发
	WSMTTaskControl      WSMsgType = 20010 // ws任务控制消息，如启动、终止、暂停，console下发
	WSMTTaskNotification WSMsgType = 20020 // ws任务通报，agent主动汇报
)

type (
	WSMessage struct {
		ID    string          `json:"id"`
		Alias string          `json:"alias,omitempty"`
		Type  WSMsgType       `json:"type"`
		From  WSMsgDirection  `json:"from,omitempty"`
		Data  json.RawMessage `json:"data,omitempty"`
	}

	// WSMConnectionRequest ws建立连接请求
	WSMConnectionRequest struct{}
	// WSMConnectionResponse ws建立连接请求的相应，发起方（client）通过Code字段来判断是否成功
	WSMConnectionResponse struct{ *BaseRes }
)
