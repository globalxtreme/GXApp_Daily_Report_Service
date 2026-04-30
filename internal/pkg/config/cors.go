package config

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/rs/cors"
)

var (
	CorsOptions cors.Options
)

func InitCors() {
	CorsOptions.AllowedOrigins = []string{xtremepkg.HostFull, "http://localhost:3000"}
	CorsOptions.AllowCredentials = true
	CorsOptions.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	CorsOptions.AllowedHeaders = []string{"*"}
}
