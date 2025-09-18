package utils

import (
	"log"
	"os"
)

func CheckEnvs() {
	if os.Getenv("IDP_BASE_URL") == "" {
		log.Fatal("IDP_BASE_URL environment variable not set")
	}
	if os.Getenv("CLIENT_ID") == "" {
		log.Fatal("CLIENT_ID environment variable not set")
	}
	if os.Getenv("CLIENT_SECRET") == "" {
		log.Fatal("CLIENT_SECRET environment variable not set")
	}
	if os.Getenv("ACCESS_REQUIREMENTS_FOR_DEVICES_File") == "" {
		log.Fatal("ACCESS_REQUIREMENTS_FOR_DEVICES_File environment variable not set")
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
