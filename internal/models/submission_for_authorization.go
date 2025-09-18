package models

type SubmissionForAuthorization struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Description string `json:"description"`
	DeviceID    string `json:"device_id"`
}

func (data *SubmissionForAuthorization) GetUserId() (string, error) {
	return data.UserID, nil
}

func (data *SubmissionForAuthorization) GetDeviceId() (string, error) {
	return data.DeviceID, nil
}

func (data *SubmissionForAuthorization) GetUserName() (string, error) {
	return data.UserName, nil
}
