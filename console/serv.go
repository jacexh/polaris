package console

import (
	"net/http"

	"github.com/jacexh/polaris/log"
	"go.uber.org/zap"
)

func Serve() {
	log.Logger.Info("start console")
	http.HandleFunc("/ws/register", register)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello, polaris!")) })
	log.Logger.Error(
		"server down",
		zap.Error(http.ListenAndServe(":16666", nil)),
	)
}
