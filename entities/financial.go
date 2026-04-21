package entities

import "time"

type AddFinancial struct {
	Id          int16
	UserId      string
	Date        time.Time `validate:"required" label:"Tanggal Transaksi"`
	Type        string    `validate:"required"`
	Nominal     int64     `validate:"required,numeric"`
	Category    string    `validate:"required" label:"Kategori"`
	Description *string   //KARENA OPTIONAL BISA DI ISI ATAU ENGGAK DATANYA JADINYA POINTER*
	Attachment  *string   //KARENA OPTIONAL BISA DI ISI ATAU ENGGAK DATANYA JADINYA POINTER*
}

type Financial struct {
	Id          int16
	UserId      string
	Date        time.Time
	Type        string
	Nominal     int64
	Category    string
	Description *string
	Attachment  *string
}
