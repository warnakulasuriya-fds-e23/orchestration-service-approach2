package models

type IncomingDataFromHikCentral struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (data *IncomingDataFromHikCentral) GetUserId() (string, error) {
	// Implement your logic to extract user ID from the incoming data
	return data.ID, nil
}
