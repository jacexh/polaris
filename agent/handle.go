package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jacexh/polaris/log"
	"github.com/jacexh/polaris/model"
	"go.uber.org/zap"
)

type (
	RequestHandle func(r *http.Request)

	// Redirector .
	Redirector struct {
		Src string // eg. test.example.com:81
		Dst string // eg. https://stage.example.com:90/asdx123
	}

	ConsolePrinter struct{}
)

func (red *Redirector) Handle(req *http.Request) {
	if red.Src != "" && req.URL.Host != red.Src {
		return
	}
	newUrl := red.Dst + req.URL.RequestURI()
	URL, err := url.Parse(newUrl)
	if err != nil {
		log.Logger.Error("error redirect request to "+red.Dst, zap.Error(err))
		return
	}
	req.URL = URL
}

func (cp ConsolePrinter) Handle(req *http.Request) {
	var r *model.Request
	r, err := model.NewFromHTTPRequest(req)
	if err != nil {
		log.Logger.Error("convert http.Request to model.Request failed", zap.Error(err))
		return
	}
	content, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Logger.Error("marshal model.Request failed", zap.Error(err))
		return
	}
	fmt.Println(string(content))
}
