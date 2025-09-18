package utils

import (
	"log"
	"os"
)

func CheckEnvs() {
	if os.Getenv("ENABLE_INTERNAL_SUBSCRIPTION_TO_FACE_MATCH_EVENT") == "" {
		log.Fatal("ENABLE_INTERNAL_SUBSCRIPTION_TO_FACE_MATCH_EVENT environment variable not set")
	}
	if os.Getenv("ENABLE_INTERNAL_SUBSCRIPTION_TO_FACE_MATCH_EVENT") == "true" && os.Getenv("ORCHESTRATION_SERVICE_BASE_URL") == "" {
		log.Fatal("ORCHESTRATION_SERVICE_BASE_URL environment variable not set")
	}
	if os.Getenv("IDP_BASE_URL") == "" {
		log.Fatal("IDP_BASE_URL environment variable not set")
	}
	if os.Getenv("IDP_CLIENT_ID") == "" {
		log.Fatal("IDP_CLIENT_ID environment variable not set")
	}
	if os.Getenv("IDP_CLIENT_SECRET") == "" {
		log.Fatal("IDP_CLIENT_SECRET environment variable not set")
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
	if os.Getenv("HCP_IP_ADDRESS") == "" {
		log.Fatal("HCP_IP_ADDRESS environment variable not set")
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
