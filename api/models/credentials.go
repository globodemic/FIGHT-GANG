package models

//Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Name     string `json:"username"`
}

//CredentialsSlice slices
type CredentialsSlice []Credentials
