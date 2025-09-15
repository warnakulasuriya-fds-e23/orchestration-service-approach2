package models

type PersonalInfoRequestObj struct {
	PersonId   string `json:"personId"`
	AppendInfo []int  `json:"appendInfo"`
}
