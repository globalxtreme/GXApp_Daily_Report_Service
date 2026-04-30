package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"service/internal/app/api/web"
	"service/internal/app/api/web/handler"
	"service/internal/pkg/config"
	"service/internal/pkg/middleware"
	reportrepository "service/internal/report/repository"
	reportservice "service/internal/report/service"
)

func Register(router *mux.Router) {
	version := os.Getenv("VERSION")
	service := os.Getenv("SERVICE")

	api := router.PathPrefix("/api").Subrouter()

	// ── Legacy web routes (gorilla/mux) ─────────────────────────────────────
	web.Register(api.PathPrefix(fmt.Sprintf("/web/%s/%s", version, service)).Subrouter())

	// ── v1 REST API ──────────────────────────────────────────────────────────
	v1 := api.PathPrefix("/v1").Subrouter()

	// Wire report services (config.PgSQL is ready by the time Register is called)
	queryRepo := reportrepository.NewQuery(config.PgSQL)
	queryService := reportservice.NewQuery(queryRepo)
	exportService := reportservice.NewExport(queryRepo)

	authHandler := handler.AuthHandler{}
	reportHandler := handler.NewReportAPIHandler(queryService, exportService)

	// Public auth routes
	v1.HandleFunc("/auth/redirect", authHandler.Redirect).Methods(http.MethodGet)
	v1.HandleFunc("/auth/callback", authHandler.Callback).Methods(http.MethodGet)

	// Protected auth routes
	v1.Handle("/auth/me", middleware.AuthMiddleware(http.HandlerFunc(authHandler.Me))).Methods(http.MethodGet)

	// Protected report routes
	v1.Handle("/reports", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.List))).Methods(http.MethodGet)
	v1.Handle("/reports/by-user", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.ListByUser))).Methods(http.MethodGet)
	v1.Handle("/reports/by-date", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.ListByDate))).Methods(http.MethodGet)
	v1.Handle("/reports/export", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.Export))).Methods(http.MethodGet)
	v1.Handle("/reports/{id:[0-9]+}", middleware.AuthMiddleware(http.HandlerFunc(reportHandler.GetByID))).Methods(http.MethodGet)
}
