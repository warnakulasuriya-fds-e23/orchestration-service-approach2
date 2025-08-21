package utils

import (
	"fmt"
)

func RoleBasedAuthorization(deviceId string, userRoles []string) (bool, error) {
	access_granted := false
	accessRequirements, err := ReadAccessRequirementsFile()
	if err != nil {
		err = fmt.Errorf("failed to read access requirements file: %w", err)
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
