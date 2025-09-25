package models

type AccessRequirementsFileStruct struct {
	Policies []Requirement `json:"policies"`
}

type Requirement struct {
	BiometricDeviceId string `json:"biometric_device_id"`
	DoorId            string `json:"door_id"`
	RequiredGroup     string `json:"required_group"`
}
