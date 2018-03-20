package agent

import (
	"errors"
	"net/url"

	"io/ioutil"

	"github.com/gorilla/websocket"
	"github.com/jacexh/polaris/log"
	"github.com/jacexh/polaris/model"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type (
	WSClient struct {
		conn *websocket.Conn
		url  *url.URL
	}
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Connect 连接指定的web socket接口
func (wc *WSClient) Connect(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	wc.conn = conn
	wc.url = u
	return nil
}

// NewWSClient 实例化WSClient
func NewWSClient(addr string) (*WSClient, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	conn, ret, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}
	log.Logger.Info("got response", zap.ByteString("body", body))
	return &WSClient{conn: conn, url: u}, nil
}

// Register 将node的信息注册的console
func (wc *WSClient) Register() error {
	req := &model.RegisterRequest{
		AgentID: uuid.NewV4().String(),
		Alias:   uuid.NewV4().String(),
	}
	data, _ := json.Marshal(req)
	err := wc.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	_, msg, err := wc.conn.ReadMessage()
	if err != nil {
		return err
	}
	log.Logger.Info("received message from server", zap.ByteString("data", msg))

	res := new(model.RegisterResponse)
	err = json.Unmarshal(msg, res)
	if err != nil {
		return err
	}

	if res.Code != model.CodeSucc {
		return errors.New(res.Message)
	}
	wc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return nil
}
