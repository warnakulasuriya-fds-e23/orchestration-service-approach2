package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
)

type RequirementsManager struct {
	AccessRequirements models.AccessRequirements
	IsInitialized      bool
	ConfiguredFilePath string
}

var (
	instance *RequirementsManager
	once     sync.Once
)

func GetRequirementsManager() *RequirementsManager {
	once.Do(func() {
		requirementFilePath, err := filepath.Abs(os.Getenv("ACCESS_REQUIREMENTS_FOR_DEVICES_File"))
		if err != nil {
			log.Fatalf("Failed to get absolute path of access requirements file: %v", err)
		}
		accessRequirements, err := ReadAccessRequirementsFile(requirementFilePath)
		if err != nil {
			log.Fatalf("Failed to read access requirements file: %v", err)
		}
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
		instance = &RequirementsManager{
			AccessRequirements: accessRequirements,
			IsInitialized:      true,
			ConfiguredFilePath: requirementFilePath,
		}
	})
	return instance
}

func (rm *RequirementsManager) GetAccessRequirements() models.AccessRequirements {
	if !rm.IsInitialized {
		log.Println("Requirements manager is not initialized")
		return models.AccessRequirements{}
	}
	return rm.AccessRequirements
}

func (rm *RequirementsManager) ReloadRequirements() error {
	accessRequirements, err := ReadAccessRequirementsFile(rm.ConfiguredFilePath)
	if err != nil {
		return fmt.Errorf("failed to reload access requirements file: %v", err)
	}
	rm.AccessRequirements = accessRequirements
	return nil
}

func (rm *RequirementsManager) GetConfiguredRequirementsFilePath() string {
	if !rm.IsInitialized {
		log.Println("Requirements manager is not initialized")
		return ""
	}
	return rm.ConfiguredFilePath
}

func (rm *RequirementsManager) GetRequiredRoleOfDevice(deviceId string) (string, error) {
	if !rm.IsInitialized {
		log.Println("Requirements manager is not initialized")
		return "", fmt.Errorf("requirements manager is not initialized")
	}
	deviceData, ok := rm.AccessRequirements.Requirements[deviceId]
	if !ok {
		log.Printf("Device ID %s not found in access requirements", deviceId)
		return "", fmt.Errorf("device ID %s not found in access requirements", deviceId)
	}
	return deviceData.RequiredRole, nil
}

func (rm *RequirementsManager) GetDoorId(deviceId string) (string, error) {
	if !rm.IsInitialized {
		log.Println("Requirements manager is not initialized")
		return "", fmt.Errorf("requirements manager is not initialized")
	}
	deviceData, ok := rm.AccessRequirements.Requirements[deviceId]
	if !ok {
		log.Printf("Device ID %s not found in access requirements", deviceId)
		return "", fmt.Errorf("device ID %s not found in access requirements", deviceId)
	}
	return deviceData.DoorId, nil
}
