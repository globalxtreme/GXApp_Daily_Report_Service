package web

import (
	"github.com/gorilla/mux"
	"service/internal/app/api/web/handler"
)

func Register(router *mux.Router) {
	activityRouter(router)
}

func activityRouter(router *mux.Router) {
	router.HandleFunc("/activities", handler.ActivityHandler{}.Get).Methods("GET")
}
