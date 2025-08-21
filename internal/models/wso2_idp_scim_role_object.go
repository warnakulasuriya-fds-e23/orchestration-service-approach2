package models

// the way that roles will be recieved through SCIM call to WSO2 IDP
type WSO2IDPRoleObject struct {
	Ref             string `json:"$ref"`
	AudienceDisplay string `json:"audienceDisplay"`
	AudienceType    string `json:"audienceType"`
	AudienceValue   string `json:"audienceValue"`
	Display         string `json:"display"`
	Value           string `json:"value"`
}

func (role *WSO2IDPRoleObject) GetRoleName() string {
	return role.Display
}
