package config

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/rs/cors"
)

var (
	CorsOptions cors.Options
)

func InitCors() {
	CorsOptions.AllowedOrigins = []string{xtremepkg.HostFull}
	CorsOptions.AllowCredentials = false
	CorsOptions.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	CorsOptions.AllowedHeaders = []string{"*"}
}
