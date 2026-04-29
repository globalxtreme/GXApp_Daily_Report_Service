package config

import "os"

var (
	PrivateHost          string
	PrivateAPIClient     map[string]map[string]interface{}
	PrivateAPICredential map[string]map[string]interface{}
)

func InitPrivateAPI() {
	PrivateHost = os.Getenv("PRIVATE_HOST")

	InitPrivateAPIClient()
	InitPrivateAPICredential()
}

func InitPrivateAPIClient() {
	PrivateAPIClient = map[string]map[string]interface{}{
		"testing": {
			"host":          os.Getenv("CLIENT_PRIVATE_API_TESTING_HOST"),
			"client-id":     os.Getenv("CLIENT_PRIVATE_API_TESTING_ID"),
			"client-name":   os.Getenv("CLIENT_PRIVATE_API_TESTING_NAME"),
			"client-secret": os.Getenv("CLIENT_PRIVATE_API_TESTING_SECRET"),
		},
	}
}

func InitPrivateAPICredential() {
	PrivateAPICredential = map[string]map[string]interface{}{
		"testing": {
			"id":  os.Getenv("PRIVATE_API_TESTING_ID"),
			"key": os.Getenv("PRIVATE_API_TESTING_KEY"),
		},
	}
}
