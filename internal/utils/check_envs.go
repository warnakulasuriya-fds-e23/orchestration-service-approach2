package utils

import (
	"log"
	"os"
)

func CheckEnvs() {
	if os.Getenv("IDP_ADDRESS") == "" {
		log.Fatal("IDP_ADDRESS environment variable not set")
	}
	if os.Getenv("IDP_USERNAME") == "" {
		log.Fatal("IDP_USERNAME environment variable not set")
	}
	if os.Getenv("IDP_PASSWORD") == "" {
		log.Fatal("IDP_PASSWORD environment variable not set")
	}
	if os.Getenv("ACCESS_REQUIREMENTS_FOR_DEVICES_File") == "" {
		log.Fatal("ACCESS_REQUIREMENTS_FOR_DEVICES_File environment variable not set")
	}
	if os.Getenv("HCP_OPENAPI_USER_KEY") == "" {
		log.Fatal("HCP_OPENAPI_USER_KEY environment variable not set")
	}
	if os.Getenv("HCP_OPENAPI_USER_SECRET") == "" {
		log.Fatal("HCP_OPENAPI_USER_SECRET environment variable not set")
	}
	if os.Getenv("HCP_ADDRESS") == "" {
		log.Fatal("HCP_ADDRESS environment variable not set")
	}
	if os.Getenv("ACCESS_CONTROL_CONFIG_INTERVAL") == "" {
		log.Fatal("ACCESS_CONTROL_CONFIG_INTERVAL environment variable not set")
	}
	if os.Getenv("ACCESS_CONTROL_CONFIG_API_KEY") == "" {
		log.Fatal("ACCESS_CONTROL_CONFIG_API_KEY environment variable not set")
	}
	if os.Getenv("ACCESS_CONTROL_CONFIG_BASE_URL") == "" {
		log.Fatal("ACCESS_CONTROL_CONFIG_BASE_URL environment variable not set")
	}
}
