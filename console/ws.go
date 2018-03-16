package console

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jacexh/polaris/log"
	"github.com/jacexh/polaris/model"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{}
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
)

func commonServerError(err error) []byte {
	data, _ := json.Marshal(&model.BaseResponse{Code: model.CodeServFail, Message: err.Error()}) // 如果你也错，我还能做什么
	return data
}

func register(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.Close()
				break
			}
			log.Logger.Error("occur error when read message", zap.Error(err), zap.Int("type", mt))
			break
		}
		log.Logger.Info("received", zap.ByteString("data", message), zap.Int("type", mt))
		// todo: save node info
		resp := &model.BaseResponse{Code: model.CodeSucc}
		var data []byte
		data, err = json.Marshal(resp)
		if err != nil {
			log.Logger.Error("occur error when marshal json object", zap.Error(err))
			data = commonServerError(err)
		}
		err = c.WriteMessage(mt, data)
		log.Logger.Info("send", zap.ByteString("data", data))
		if err != nil {
			log.Logger.Error("occur error when write message", zap.Error(err))
			break
		}
	}
}
