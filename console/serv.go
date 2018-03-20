package console

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jacexh/polaris/log"
	"go.uber.org/zap"
)

func Serve() {
	log.Logger.Info("start console")
	router := mux.NewRouter()
	router.HandleFunc("/api/register", register).Methods(http.MethodPost)
	router.HandleFunc("/ws", ws)

	log.Logger.Error("server down", zap.Error(http.ListenAndServe(":16666", router)))
}
