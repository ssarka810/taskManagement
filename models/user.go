package models

type UserDetails struct{
	UserName string `json:"username,omitempty" validate:"required,min=4"`
	Password string `json:"password,omitempty" validate:"required,min=4"`
}