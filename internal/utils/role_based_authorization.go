package utils

import (
	"fmt"
)

func RoleBasedAuthorization(deviceId string, userRoles []string) (bool, error) {
	access_granted := false
	requirementsManager := GetRequirementsManager()
	accessRequirements := requirementsManager.GetAccessRequirements()
	if !requirementsManager.IsInitialized {
		err := fmt.Errorf("requirementsManager is not initialized")
		return access_granted, err
	}
	if accessRequirements.Requirements == nil {
		err := fmt.Errorf("no access requirements found, please check the file pointed to in environmental variable ACCESS_REQUIREMENTS_FILE")
		return access_granted, err
	}
	DeviceData, ok := accessRequirements.Requirements[deviceId]
	requiredRole := DeviceData.RequiredRole
	if !ok {
		err := fmt.Errorf("device ID %s not found in access requirements", deviceId)
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
