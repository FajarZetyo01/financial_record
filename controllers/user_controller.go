package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	db *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{db: db}
}

func (controller UserController) Profile(writer http.ResponseWriter, request *http.Request) {
	layout := "views/user/profile.html"

	data := make(map[string]interface{})

	//MENAMPILKAN FLASH MESSAGE DARI LOGIN
	flash := config.SessionManager.PopString(request.Context(), "success")
	if flash != "" {
		data["success"] = flash
	}

	//AMBIL USER_ID DARI SESSION
	sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")

	//TAMPILKAN DATA USER
	userData, err := models.NewUserModel(controller.db).FindUserById(sessionUserId)
	if err != nil {
		data["error"] = "User tidak ditemukan" + err.Error()
	} else {
		data["user"] = userData
	}
	if request.Method == http.MethodPost {
		//BATASI UKURAN FOTO SEBESAR 5MB
		request.ParseMultipartForm(5 * 1024 * 1024)

		//AMBIL DATA YANG MAU DIUBAH
		password := request.FormValue("password")
		user := entities.User{
			Id:       sessionUserId,
			Password: password,
			Name:     request.FormValue("name"),
			Photo:    userData.Photo,
		}

		//TAMPILKAN ERROR VALIDASI
		if err := helpers.NewValidation(controller.db).ValidateStruct(user); err != nil {
			data["validation"] = err
			data["user"] = userData
			views.RenderTemplate(writer, layout, data)
			return
		}
		//HASH PASSWORD BARU
		if password != "" {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			user.Password = string(hashPassword) //MASUKAN KEMBALI KE STRUCT
		}

		if file, handler, err := request.FormFile("photo"); err == nil {
			defer file.Close() //JEDA PROSES SAMPAI DIPILIH FILENYA

			//VALIDASI UKURAN
			if handler.Size > 5*1024*1024 {
				data["error"] = "Ukuran file terlalu besar, maksimal 5MB"
				data["user"] = userData
				views.RenderTemplate(writer, layout, data)
				return
			}

			//VALIDASI EXTENSI FILE
			ext := strings.ToLower(filepath.Ext(handler.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
				data["error"] = "Tipe file tidak mendukung"
				data["user"] = userData
				views.RenderTemplate(writer, layout, data)
				return
			}

			//NAMA & PATH FILE BARU
			filename := fmt.Sprintf("profile_%s%s", time.Now().Format("2006-01-02_15-04-05"), ext)
			path := filepath.Join("public/user_photo", filename)

			//HAPUS FOTO LAMA JIKA ADA
			if userData.Photo != nil && *userData.Photo != "" {
				oldPath := filepath.Join("public/user_photo", *userData.Photo)
				os.Remove(oldPath)
			}

			//BUAT FILE FOTO BARU
			out, err := os.Create(path)
			if err != nil {
				data["error"] = "Gagal membuat file baru" + err.Error()
				data["user"] = userData
				views.RenderTemplate(writer, layout, data)
				return
			}
			defer out.Close() //JEDA PROSES SAMPAI FILE BARU TERBUAT

			//SIMPAN FILE BARU
			_, errCopy := io.Copy(out, file)
			if errCopy != nil {
				data["error"] = "Gagal menyimpan file baru" + err.Error()
				data["user"] = userData
				views.RenderTemplate(writer, layout, data)
				return
			}
			//SET KE STRUCT
			user.Photo = &filename
		}

		//UPDATE PROFILE
		if err := models.NewUserModel(controller.db).UpdateProfile(user); err != nil {
			data["error"] = "Gagal mengubah data profile" + err.Error()
		} else {
			config.SessionManager.Put(request.Context(), "success", "Berhasil mengubah data profile")
			http.Redirect(writer, request, "/profile", http.StatusSeeOther)
			return
		}
	}
	views.RenderTemplate(writer, layout, data)
}
