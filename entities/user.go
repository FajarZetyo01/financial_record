package entities

type User struct {
	Id       string
	Name     string `validate:"required" label:"Nama"`
	Email    string
	Password string `validate:"omitempty,min=6"` //omitempty JIKA DI ISI MINIMAL 6
	Photo    *string
}
