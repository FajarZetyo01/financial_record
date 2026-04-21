package entities

type Register struct {
	Id              string
	Name            string `validate:"required" label:"Nama"`
	Email           string `validate:"required,email,isunique=users-email"`                         //ISUNIQUE NGECEK NAMA TABEL & KOLOMNYA
	Password        string `validate:"required,min=6"`                                              //Minimal 6 huruf & angka
	ConfirmPassword string `validate:"required,min=6,eqfield=Password" label:"Konfirmasi Password"` //eqfield isiannya harus sama dengan variabel password
}
type Auth struct {
	Id       string
	Name     string
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}
