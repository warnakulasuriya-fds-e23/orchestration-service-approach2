package utils

import (
	"log"
	"sync"

	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
)

type RequirementsManager struct {
	AccessRequirements models.AccessRequirements
	IsInitialized      bool
}

var (
	instance *RequirementsManager
	once     sync.Once
)

func GetRequirementsManager() *RequirementsManager {
	once.Do(func() {
		accessRequirements, err := ReadAccessRequirementsFile()
		if err != nil {
			log.Fatalf("Failed to read access requirements file: %v", err)
		}
		instance = &RequirementsManager{
			AccessRequirements: accessRequirements,
			IsInitialized:      true,
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
