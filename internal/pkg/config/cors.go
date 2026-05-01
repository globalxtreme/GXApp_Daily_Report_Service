package config

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/rs/cors"
	"os"
)

var (
	CorsOptions cors.Options
)

func InitCors() {
	CorsOptions.AllowedOrigins = []string{xtremepkg.HostFull, os.Getenv("FRONTEND_HOST")}
	CorsOptions.AllowCredentials = true
	CorsOptions.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	CorsOptions.AllowedHeaders = []string{"*"}
}
