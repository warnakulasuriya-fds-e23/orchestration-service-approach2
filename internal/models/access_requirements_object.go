package models

//TODO: finalize struct to which the access requirement data will be loaded

type AccessRequirements struct {
	Requirements map[string]DeviceData `json:"requirements"`
}

type DeviceData struct {
	DoorId       string `json:"doorId"`
	RequiredRole string `json:"requiredRole"`
}

// TODO: finalize get method to get the role with necessary error handling
func (ar *AccessRequirements) GetRequiredRole(deviceId string) (role string, err error) {
	deviceData := ar.Requirements[deviceId]
	return deviceData.RequiredRole, nil
}
func (ar *AccessRequirements) GetDoorId(deviceId string) (doorId string, err error) {
	deviceData := ar.Requirements[deviceId]
	return deviceData.DoorId, nil
}
