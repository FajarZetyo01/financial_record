package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *sql.DB //dependency database Disimpan di struct supaya bisa dipakai di semua method controller ada fo
	//folder config/database.go
	//Intinya biar bisa akses database nya
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		db: db,
	}
}

/*
return (mengembalikan ) &AuthController{ ... }
Artinya:
Buat object AuthController
Ambil alamat memorinya alamat pointernya adalah *AuthController
Return pointer
➡️ Karena function-nya return *AuthController

db: db
Posisi	Artinya
kiri  db ->nama field struct disini adalah db *sql.DB
kanan db ->variable parameter function
Kenapa namanya sama? Bikin bingung 😅
Karena:
konvensi Go
field dan parameter sering dinamai sama
*/

//📌 Tujuan:
//Membuat instance AuthController
//Menerapkan Dependency Injection
//Biar AuthController selalu punya akses ke DB
//CONTOH PEMAKAIANYA DI ROUTE : authController := NewAuthController(db)
/*
	NewAuthController(db)
          ↑
   db berasal dari sql.Open()
          ↑
   database/sql package
          ↑
   Database asli (MySQL / PostgreSQL / dll)
*/

// func controller authcontroller ini jadinya method karena mengiket ke struct
func (controller *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
	//Kenapa ada pointer *AuthController
	//Kenapa?
	//Konsisten dengan NewAuthController() yang return *AuthController
	//Lebih hemat memory
	//Aman kalau nanti ada state di controller

	//file layout
	template := "views/auth/register.html" //template → file HTML

	//KIRIM DATA KE HTML
	data := make(map[string]interface{}) //data → data yang dikirim ke HTML (validation, error, success, dll)

	//KETIKA BUTTON REGISTER DI KLIK
	if request.Method == http.MethodPost {
		request.ParseForm()
		register := entities.Register{
			Id:              uuid.New().String(),       //HASILNYA DI DB JADI UNIQ IDNYA
			Name:            request.FormValue("name"), //value name di ambil dari name di register.html
			Email:           request.FormValue("email"),
			Password:        request.FormValue("password"),
			ConfirmPassword: request.FormValue("confirm_password"),
		}
		//📌 Yang terjadi:
		//Ambil data dari file html <input name="...">
		//Simpan ke struct entities.Register
		//uuid.New().String() → bikin ID unik

		//TAMPILKAN ERROR DARI VALIDATOR
		if err := helpers.NewValidation(controller.db).ValidateStruct(register); err != nil {
			data["validation"] = err
			data["register"] = register
			views.RenderTemplate(writer, template, data)
			return //RETURN MENGEMBALIKAN KE BROWSER VIEWS RENDER TEMPLATE
		}
		//📌 Alurnya:
		//Jalankan validator
		//Kalau error:
		//kirim error ke view
		//kirim ulang data input (biar form tidak kosong)
		//render ulang halaman register
		//➡️ return penting supaya tidak lanjut ke proses berikutnya

		//HASH PASSWORD
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(writer, "Failed to secure password", http.StatusInternalServerError)
			return
		}
		register.Password = string(hashedPassword)

		//📌 Tujuan:
		//Password tidak disimpan plaintext
		//bcrypt → hashing aman
		//DefaultCost → tingkat keamanan standar

		//INSERT KE DB
		if err := models.NewAuthModel(controller.db).Register(register); err != nil {
			data["error"] = "Registrasi Gagal" + err.Error()
			views.RenderTemplate(writer, template, data)
			return
			//📌 Artinya:
			//Panggil model
			//Simpan user ke DB
			//Kalau gagal → tampilkan pesan error
		} else {
			data["success"] = "Registrasi berhasil, silahkan login"
			views.RenderTemplate(writer, template, data)
			return
			//http.Redirect(writer, request, "/login", http.StatusSeeOther)
			//return
			//📌 Alur:
			//Set pesan sukses
			//Redirect ke /login
			//StatusSeeOther → standar redirect setelah POST (anti double submit)
		}
	}
	views.RenderTemplate(writer, template, nil)
	//views.RenderTemplate(writer, template, nil) -> 1️⃣1️⃣ Request GET (Pertama kali buka halaman)
	//FLOW ALL
	//GET /register
	//→ tampilkan form
	//POST /register
	//→ ambil input
	//→ validasi
	//→ hash password
	//→ simpan ke DB
	//→ redirect ke login
}
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	template := "views/auth/login.html"

	//KIRIM DATA KE HTML
	data := make(map[string]interface{})

	//AMBIL INPUTAN DARI HTML
	if r.Method == http.MethodPost {
		r.ParseForm()
		login := entities.Auth{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		// TAMPILKAN ERROR DARI VALIDATOR
		if err := helpers.NewValidation(c.db).ValidateStruct(login); err != nil {
			data["validation"] = err
			data["login"] = login
			views.RenderTemplate(w, template, data)
			return
		}
		//CARI USER BERDASARKAN EMAIL
		user, err := models.NewAuthModel(c.db).Login(login.Email)
		if err != nil {
			data["error"] = "Akun tidak ditemukan " + err.Error()
			views.RenderTemplate(w, template, data)
			return
		}
		//COCOKAN PASSWORD
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
			data["error"] = "Password salah!"
			views.RenderTemplate(w, template, data)
			return
		}
		//SIMPAN DATA KE SESSION
		//LOGIN, USER_ID
		config.SessionManager.Put(r.Context(), "LOGGED_IN", true)
		config.SessionManager.Put(r.Context(), "USER_ID", user.Id)

		//FLASH MESSAGE
		config.SessionManager.Put(r.Context(), "success", "Selamat Datang "+user.Name)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	views.RenderTemplate(w, template, nil)
}

func (controller AuthController) Logout(writer http.ResponseWriter, request *http.Request) {
	//HAPUS SEMUA SESSION
	err := config.SessionManager.Destroy(request.Context())
	if err != nil {
		http.Error(writer, "Gagal Logout", http.StatusSeeOther)
		return
	}
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}
