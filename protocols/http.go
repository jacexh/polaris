package protocols

const (
	CodeSucc       ReturnCode = 200 // 业务处理成功
	CodeClientFail ReturnCode = 400 // 业务处理失败，client传入参数存在问题
	CodeServFail   ReturnCode = 500 // 业务处理失败，server端处理异常
)

type (
	RegisterRequest struct {
		*AgentNode
		IPs []string `json:"ips"`
	}

	RegisterResponse struct {
		*BaseRes
	}
)
