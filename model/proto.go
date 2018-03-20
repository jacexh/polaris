package model

type (
	Code int

	BaseResponse struct {
		Code    Code   `json:"code"`              // 业务处理状态码
		Message string `json:"message,omitempty"` // 业务处理成功时，该字段不一定返回；业务处理失败，该字段必定会返回
	}

	Node struct {
		ID    string `json:"id"`
		Alias string `json:"alias,omitempty"`
	}

	// RegisterRequest /api/register 请求报文
	RegisterRequest struct {
		*Node
		IPs []string `json:"ips"`
	}

	// RegisterResponse /api/register 服务端相应
	RegisterResponse struct {
		*BaseResponse
	}

	// WSConnect web socket client连接server时的请求报文
	WSConnect struct {
		*Node
	}
)

const (
	CodeSucc       Code = 200 // 业务处理成功
	CodeClientFail Code = 400 // 业务处理失败，client传入参数存在问题
	CodeServFail   Code = 500 // 业务处理失败，server端处理异常
)
