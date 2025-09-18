package models

type AccessRequirementsFileStruct struct {
	Requirements []Requirement `json:"requirements"`
}

type Requirement struct {
	HCPDefinedCameraName string `json:"hcp_defined_camera_name"`
	DoorId               string `json:"door_id"`
	RequiredRole         string `json:"required_role"`
}
