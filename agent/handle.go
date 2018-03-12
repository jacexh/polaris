package agent

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jacexh/polaris/log"
	"github.com/jacexh/polaris/model"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type (
	// RequestHandle 声明处理*http.Request对象的方法类型
	RequestHandle func(r *http.Request)

	// Redirector 实现将*http.Request转发到目标服务器
	Redirector struct {
		Src  string // eg. test.example.com:81
		Dst  string // eg. https://stage.example.com:90/asdx123
		Host string // 非空字符串时，会修改Header的Host字段
	}

	// ConsolePrinter 将*http.Request打印到控制台
	ConsolePrinter struct{}
)

// Handle *Redirector实现的RequestHandle
func (red *Redirector) Handle(req *http.Request) {
	if red.Src != "" && req.URL.Host != red.Src {
		return
	}
	newURL := red.Dst + req.URL.RequestURI()
	URL, err := url.Parse(newURL)
	if err != nil {
		log.Logger.Error("error redirect request to "+red.Dst, zap.Error(err))
		return
	}
	req.URL = URL
	if red.Host != "" {
		req.Header.Set("Host", red.Host)
	}
}

// Handle ConsolePrinter实现的RequestHandle
func (cp ConsolePrinter) Handle(req *http.Request) {
	var r *model.Request
	r, err := model.NewFromHTTPRequest(req)
	if err != nil {
		log.Logger.Error("convert http.Request to model.Request failed", zap.Error(err))
		return
	}
	content, err := jsoniter.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Logger.Error("marshal model.Request failed", zap.Error(err))
		return
	}
	fmt.Println(string(content))
}
