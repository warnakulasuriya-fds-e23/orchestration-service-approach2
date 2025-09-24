package utils

import (
	"fmt"
	"strings"
)

func GroupBasedAuthorization(deviceId string, userGroups []string) (bool, error) {
	access_granted := false
	requirementsManager := GetRequirementsManager()
	if !requirementsManager.IsInitialized {
		err := fmt.Errorf("requirementsManager is not initialized")
		return access_granted, err
	}

	requiredGroup, err := requirementsManager.GetRequiredGroupOfDevice(deviceId)
	if err != nil {
		err := fmt.Errorf("required Group of device with device id %s couldn't be internally accessed: %v", deviceId, err)
		return access_granted, err
	}
	for _, userGroup := range userGroups {
		if strings.Contains(userGroup, requiredGroup) {
			access_granted = true
			return access_granted, nil
		}
	}
	return access_granted, nil
}
