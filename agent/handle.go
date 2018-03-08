package agent

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin/json"
	"github.com/jacexh/polaris/log"
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
	header, err := json.Marshal(req.Header)
	if err != nil {
		log.Logger.Error("error marshal request header", zap.Error(err))
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.Error("read request body failed", zap.Error(err))
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	fmt.Printf("headers:\n%s\n\nbody:\n%s\n", header, body)
}
