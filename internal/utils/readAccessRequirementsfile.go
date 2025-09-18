package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
)

//TODO: Read accessRequirementsForDevices.json and get the data loaded in to an
//access requirements object

func ReadAccessRequirementsFile(requirementFilePath string) (models.AccessRequirements, error) {
	var accessRequirementsFileStruct models.AccessRequirementsFileStruct
	var accessRequirements models.AccessRequirements
	reqFile := os.Getenv("ACCESS_REQUIREMENTS_FOR_DEVICES_File")
	reqFilePath, err := filepath.Abs(reqFile)
	if err != nil {
		return accessRequirements, fmt.Errorf("failed to get absolute path of the access requirements for devices file: %w", err)
	}
	file, err := os.Open(reqFilePath)
	if err != nil {
		err = fmt.Errorf("failed to open access requirements file: %w", err)
		return accessRequirements, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&accessRequirementsFileStruct)
	if err != nil {
		return accessRequirements, err
	}

	// Map the requirements from the file struct to the access requirements struct
	accessRequirements.Requirements = make(map[string]models.DeviceData)
	for _, requirement := range accessRequirementsFileStruct.Requirements {
		accessRequirements.Requirements[requirement.BiometricDeviceId] = models.DeviceData{
			DoorId:       requirement.DoorId,
			RequiredRole: requirement.RequiredRole,
		}
	}

	return accessRequirements, nil
}
