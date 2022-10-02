package models

type User struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	Email          string `json:"mail"`
	JobTitle       string `json:"jobTitle"`
	OfficeLocation string `json:"officeLocation"`
	Token          string `json:"token"`
}
