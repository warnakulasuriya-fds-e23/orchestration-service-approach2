package models

type AccessRequirementsFileStruct struct {
	Requirements []Requirement `json:"requirements"`
}

type Requirement struct {
	BiometricDeviceId string `json:"biometric_device_id"`
	DoorId            string `json:"door_id"`
	RequiredRole      string `json:"required_role"`
}
