package model

type (
	Code int

	BaseResponse struct {
		Code    Code   `json:"code"`              // 业务处理状态码
		Message string `json:"message,omitempty"` // 业务处理成功时，该字段不一定返回；业务处理失败，该字段必定会返回
	}

	RegisterRequest struct {
		AgentID string   `json:"id"`
		Alias   string   `json:"alias"`
		IPs     []string `json:"ips"`
	}

	RegisterResponse struct {
		*BaseResponse
	}
)

const (
	CodeSucc       Code = 200 // 业务处理成功
	CodeClientFail Code = 400 // 业务处理失败，client传入参数存在问题
	CodeServFail   Code = 500 // 业务处理失败，server端处理异常
)
