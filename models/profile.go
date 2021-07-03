package models

type UserProfile struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
