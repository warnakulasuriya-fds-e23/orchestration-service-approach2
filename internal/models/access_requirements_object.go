package models

//TODO: finalize struct to which the access requirement data will be loaded

type AccessRequirements struct {
	Requirements map[string]DeviceData `json:"requirements"`
}

type DeviceData struct {
	DoorId        string `json:"doorId"`
	RequiredGroup string `json:"requiredGroup"`
}

// TODO: finalize get method to get the group with necessary error handling
func (ar *AccessRequirements) GetRequiredGroup(deviceId string) (group string, err error) {
	deviceData := ar.Requirements[deviceId]
	return deviceData.RequiredGroup, nil
}
func (ar *AccessRequirements) GetDoorId(deviceId string) (doorId string, err error) {
	deviceData := ar.Requirements[deviceId]
	return deviceData.DoorId, nil
}
