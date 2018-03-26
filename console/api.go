package console

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jacexh/polaris/console/model"
	"github.com/jacexh/polaris/log"
	"github.com/jacexh/polaris/protocols"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	upgrader = websocket.Upgrader{}
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
)

type (
	wsServer struct {
		id   string
		conn *websocket.Conn
		out  chan *message
	}

	connPool sync.Map

	message struct {
		t    int
		data []byte
	}
)

func newWSServer(w http.ResponseWriter, r *http.Request) (*wsServer, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn.SetPingHandler(nil)
	conn.SetCloseHandler(nil)
	conn.SetPongHandler(nil)

	// todo: 解析agent id
	req := new(protocols.WSMessage)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}
	return &wsServer{
		id:   req.ID,
		conn: conn,
		out:  make(chan *message),
	}, nil
}

func (ws *wsServer) readPump() {
	defer func() {
		close(ws.out)
		ws.conn.Close()
	}()
	for {
		mt, msg, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Logger.Error("connection lost", zap.Error(err), zap.String("node", ws.id))
			}
			log.Logger.Warn("occur error when read message", zap.String("node", ws.id), zap.Error(err))
		}

		switch mt {
		case websocket.TextMessage:
			// todo: 处理业务消息
			ret := new(protocols.WSMessage)
			err = json.Unmarshal(msg, ret)
			if err != nil {
				log.Logger.Warn("unmarshal message failed", zap.String("node", ws.id), zap.Error(err))
			}
		default:
			log.Logger.Warn("unsupported websocket message type", zap.Int("type", mt), zap.ByteString("data", msg))
		}
	}
}

func (ws *wsServer) writePump() {
	defer func() {
		close(ws.out)
		ws.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-ws.out:
			if !ok { // 写通道被关闭
				ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			ws.conn.WriteMessage(msg.t, msg.data)
		}
	}
}

func commonServerError(err error) []byte {
	data, _ := json.Marshal(&protocols.BaseRes{Code: protocols.CodeServFail, Message: err.Error()}) // 如果这也错，我还能做什么
	return data
}

func register(w http.ResponseWriter, r *http.Request) {
	req := new(protocols.RegisterRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		json.NewEncoder(w).Encode(&protocols.RegisterResponse{&protocols.BaseRes{Code: protocols.CodeClientFail, Message: err.Error()}})
		log.Logger.Error("unmarshal request body failed", zap.Error(err))
		return
	}
	log.Logger.Info("received", zap.Any("request", req))

	// todo: save node info

	res := &protocols.RegisterResponse{BaseRes: &protocols.BaseRes{Code: protocols.CodeSucc}}
	data, err := json.Marshal(res)
	if err != nil {
		data = commonServerError(err)
		log.Logger.Error("marshal response body failed", zap.Error(err))
	}
	w.Write(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func ws(w http.ResponseWriter, r *http.Request) {
	serv, err := newWSServer(w, r)
	if err != nil {
		log.Logger.Error("init wsServ failed", zap.Error(err))
	}
	go serv.readPump()
	serv.writePump()
}

func newAgent() *model.Agent {
	return nil
}
