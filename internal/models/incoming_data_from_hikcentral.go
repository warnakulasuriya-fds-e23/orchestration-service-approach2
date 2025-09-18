package models

type IncomingDataFromHikCentral struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Description string `json:"description"`
	DeviceID    string `json:"device_id"`
}

func (data *IncomingDataFromHikCentral) GetUserId() (string, error) {
	return data.UserID, nil
}

func (data *IncomingDataFromHikCentral) GetDeviceId() (string, error) {
	return data.DeviceID, nil
}

func (data *IncomingDataFromHikCentral) GetUserName() (string, error) {
	return data.UserName, nil
}
