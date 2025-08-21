package utils

import (
	"fmt"
)

func RoleBasedAuthorization(deviceId string, userRoles []string) (bool, error) {
	access_granted := false
	requirementsManager := GetRequirementsManager()
	if !requirementsManager.IsInitialized {
		err := fmt.Errorf("requirementsManager is not initialized")
		return access_granted, err
	}

	requiredRole, err := requirementsManager.GetRequiredRoleOfDevice(deviceId)
	if err != nil {
		err := fmt.Errorf("required role of device with device id %s couldn't be internally accessed: %v", deviceId, err)
		return access_granted, err
	}
	for _, userRole := range userRoles {
		if userRole == requiredRole {
			access_granted = true
			return access_granted, nil
		}
	}
	return access_granted, nil
}
