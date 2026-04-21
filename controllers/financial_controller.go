package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type FinancialController struct {
	db *sql.DB
}

func NewFinancialController(db *sql.DB) *FinancialController {
	return &FinancialController{
		db: db,
	}
}

// FORMAT RUPAIH
func formatIDR(n int64) string {
	str := fmt.Sprintf("%d", n)
	var result []string
	for len(str) > 3 {
		result = append([]string{str[len(str)-3:]}, result...)
		str = str[:len(str)-3]
	}
	if len(str) > 0 {
		result = append([]string{str}, result...)
	}
	return strings.Join(result, ".") + ",00"
}
func (controller FinancialController) Home(writer http.ResponseWriter, request *http.Request) {

	layout := "views/financial/home.html"
	var data = make(map[string]interface{})

	//MENAMPILKAN FLASH MESSAGE DARI LOGIN
	flash := config.SessionManager.PopString(request.Context(), "success")
	if flash != "" {
		data["success"] = flash
	}

	//MENAMPILKAN DROPDOWN BULAN
	currentDate := time.Now()
	var months []string
	for i := 0; i < 6; i++ {
		previousMonth := currentDate.AddDate(0, -i, 0)
		months = append(months, previousMonth.Format("January 2006"))
	}
	data["months"] = months //tampilkan di HTML

	//PILIH BULAN
	selectedMonth := request.URL.Query().Get("selected_month")
	if selectedMonth == "" {
		selectedMonth = currentDate.Format("January 2006")
	}
	data["selectedMonth"] = selectedMonth

	//PILIH TIPE (PEMASUKAN/PENGELUARAN)
	pemasukanOnly := request.URL.Query().Get("pemasukanOnly") == "true"
	data["pemasukanOnly"] = pemasukanOnly
	pengeluaranOnly := request.URL.Query().Get("pengeluaranOnly") == "true"
	data["pengeluaranOnly"] = pengeluaranOnly

	//AMBIL USER_ID DARI SESSION
	sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")

	//TAMPILKAN TOTAL PEMASUKAN, TOTAL PENGELUARAN
	pemasukan, pengeluaran, err := models.NewFinancialModel(controller.db).GetFinancialTotalNominal(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		data["error"] = "Gagal menampilkan total pemasukan & pengeluaran " + err.Error()
	} else {
		data["total_pemasukan"] = pemasukan
		data["total_pengeluaran"] = pengeluaran
	}

	//TAMPILKAN DATA LIST
	listFinancial, err := models.NewFinancialModel(controller.db).FindAllFinancial(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		data["error"] = "Gagal menampilkan list keuangan " + err.Error()
	} else {
		data["financials"] = listFinancial
	}

	//FUNC IDR RUPIAH
	funcMap := template.FuncMap{
		"formatIDR": formatIDR,
		"indexNo":   func(a, b int) int { return a + b },
	}
	tmpl, _ := template.New(filepath.Base(layout)).Funcs(funcMap).ParseFiles(layout)
	tmpl.Execute(writer, data)
}

// ADD FINANCIAL BUAT TESTING SEBELUM INSERT KE DB BISA PAKE CARA INI
//func (controller FinancialController) AddFinancialRecord(writer http.ResponseWriter, request *http.Request) {
//	layout := "views/financial/create.html"
//
//	var data = make(map[string]interface{})
//
//	if request.Method == http.MethodPost {
//		request.ParseForm()
//
//		//AMBIL TANGGAL
//		dateStr := request.FormValue("date")
//		date, _ := time.Parse("2006-01-02", dateStr)
//
//		//AMBIL NOMINAL
//		nominalStr := request.FormValue("nominal")
//		nominal, _ := strconv.ParseInt(nominalStr, 10, 64)
//
//		//AMBIL ATTACHMENT
//		var attachment *string
//		if attachmentValue := request.FormValue("attachment"); attachmentValue != "" {
//			attachment = &attachmentValue
//		}
//
//		//AMBIL DESKRIPSI
//		var descirption *string
//		if descriptionValue := request.FormValue("description"); descriptionValue != "" {
//			descirption = &descriptionValue
//		}
//
//		//AMBIL USER_ID DARI SESSION
//		sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")
//
//		//MASUKAN KE STRUCT
//		financial := entities.AddFinancial{
//			UserId:      sessionUserId,
//			Date:        date,
//			Type:        request.FormValue("type"),
//			Category:    request.FormValue("category"),
//			Nominal:     nominal,
//			Attachment:  attachment,
//			Description: descirption,
//		}
//
//		//TAMPILKAN ERROR
//		if err := helpers.NewValidation(controller.db).ValidateStruct(financial); err != nil {
//			data["validation"] = err
//			data["financial"] = financial
//			views.RenderTemplate(writer, layout, data)
//			return
//		}
//		fmt.Println(financial)
//	}
//	views.RenderTemplate(writer, layout, nil)
//}

// ADD FINANCIAL
func (controller FinancialController) AddFinancialRecord(writer http.ResponseWriter, request *http.Request) {
	layout := "views/financial/create.html"

	var data = make(map[string]interface{})

	if request.Method == http.MethodPost {
		request.ParseForm()

		//AMBIL TANGGAL
		dateStr := request.FormValue("date")
		date, _ := time.Parse("2006-01-02", dateStr)

		//AMBIL NOMINAL
		nominalStr := request.FormValue("nominal")
		nominal, _ := strconv.ParseInt(nominalStr, 10, 64)

		//AMBIL ATTACHMENT
		var attachment *string
		if attachmentValue := request.FormValue("attachment"); attachmentValue != "" {
			attachment = &attachmentValue
		}

		//AMBIL DESKRIPSI
		var descirption *string
		if descriptionValue := request.FormValue("description"); descriptionValue != "" {
			descirption = &descriptionValue
		}

		//AMBIL USER_ID DARI SESSION
		sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")

		//MASUKAN KE STRUCT
		financial := entities.AddFinancial{
			UserId:      sessionUserId,
			Date:        date,
			Type:        request.FormValue("type"),
			Category:    request.FormValue("category"),
			Nominal:     nominal,
			Attachment:  attachment,
			Description: descirption,
		}

		//TAMPILKAN ERROR
		if err := helpers.NewValidation(controller.db).ValidateStruct(financial); err != nil {
			data["validation"] = err
			data["financial"] = financial
			views.RenderTemplate(writer, layout, data)
			return
		}
		//fmt.Println(financial)

		//INSERT KE DATABASE
		if err := models.NewFinancialModel(controller.db).AddFinancialRecord(financial); err != nil {
			data["error"] = err
			views.RenderTemplate(writer, layout, data)
			return
		} else {
			//FLASH MESSAGE
			config.SessionManager.Put(request.Context(), "success", "Berhasil menambahkan data keuangan")
			http.Redirect(writer, request, "/home", http.StatusSeeOther)
		}
	}
	views.RenderTemplate(writer, layout, nil)
}

func (controller FinancialController) EditFinancialRecord(writer http.ResponseWriter, request *http.Request) {
	layout := "views/financial/edit.html"
	var data = make(map[string]interface{})

	//AMBIL & CHECK ID DARI URL
	idStr := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 16)
	if idStr == "" || err != nil {
		data["error"] = "ID tidak valid" + err.Error()
		views.RenderTemplate(writer, layout, data)
		return
	}

	//TAMPILKAN DATA BERDASARKAN ID
	findFinancial, err := models.NewFinancialModel(controller.db).FindFinancialById(int16(id))
	if err != nil {
		data["error"] = "Data keuangan tidak ditemukan" + err.Error()
	} else {
		data["financial"] = findFinancial
	}

	if request.Method == http.MethodPost {
		request.ParseForm()

		//AMBIL TANGGAL
		dateStr := request.FormValue("date")
		date, _ := time.Parse("2006-01-02", dateStr)

		//AMBIL NOMINAL
		nominalStr := request.FormValue("nominal")
		nominal, _ := strconv.ParseInt(nominalStr, 10, 64)

		//AMBIL ATTACHMENT
		var attachment *string
		if attachmentValue := request.FormValue("attachment"); attachmentValue != "" {
			attachment = &attachmentValue
		}

		//AMBIL DESKRIPSI
		var descirption *string
		if descriptionValue := request.FormValue("description"); descriptionValue != "" {
			descirption = &descriptionValue
		}

		//AMBIL USER_ID DARI SESSION
		sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")

		//MASUKAN KE STRUCT
		financial := entities.AddFinancial{
			Id:          int16(id),
			UserId:      sessionUserId,
			Date:        date,
			Type:        request.FormValue("type"),
			Category:    request.FormValue("category"),
			Nominal:     nominal,
			Attachment:  attachment,
			Description: descirption,
		}

		//TAMPILKAN ERROR
		if err := helpers.NewValidation(controller.db).ValidateStruct(financial); err != nil {
			data["validation"] = err
			data["financial"] = financial
			views.RenderTemplate(writer, layout, data)
			return
		}

		//UBAH DATA
		if err := models.NewFinancialModel(controller.db).EditFinancialRecord(financial); err != nil {
			data["error"] = err
			views.RenderTemplate(writer, layout, data)
			return
		} else {
			//FLASH MESSAGE
			config.SessionManager.Put(request.Context(), "success", "Berhasil mengubah data keuangan")
			http.Redirect(writer, request, "/home", http.StatusSeeOther)
		}
	}
	views.RenderTemplate(writer, layout, data)
}

func (controller FinancialController) DeleteFinancialRecord(writer http.ResponseWriter, request *http.Request) {
	//AMBIL & CHECK ID DARI URL
	idStr := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 16)
	if idStr == "" || err != nil {
		config.SessionManager.Put(request.Context(), "error", "Gagal mengambil data keuangan"+err.Error())
		http.Redirect(writer, request, "/home", http.StatusSeeOther)
		return
	}

	//HAPUS DATA
	if err := models.NewFinancialModel(controller.db).DeleteFinancialRecord(int16(id)); err != nil {
		config.SessionManager.Put(request.Context(), "error", "Gagal menghapus data keuangan"+err.Error())
	} else {
		config.SessionManager.Put(request.Context(), "success", "Berhasil menghapus data keuangan")
	}
	http.Redirect(writer, request, "/home", http.StatusSeeOther)
}

func (controller FinancialController) DownloadFinancialRecord(writer http.ResponseWriter, request *http.Request) {

	layout := "views/financial/download.html"
	var data = make(map[string]interface{})

	//PILIH BULAN
	selectedMonth := request.URL.Query().Get("selected_month")
	data["selectedMonth"] = selectedMonth

	//PILIH TIPE (PEMASUKAN/PENGELUARAN)
	pemasukanOnly := request.URL.Query().Get("pemasukanOnly") == "true"
	data["pemasukanOnly"] = pemasukanOnly
	pengeluaranOnly := request.URL.Query().Get("pengeluaranOnly") == "true"
	data["pengeluaranOnly"] = pengeluaranOnly

	//AMBIL USER_ID DARI SESSION
	sessionUserId := config.SessionManager.GetString(request.Context(), "USER_ID")

	//TAMPILKAN TOTAL PEMASUKAN, TOTAL PENGELUARAN
	pemasukan, pengeluaran, err := models.NewFinancialModel(controller.db).GetFinancialTotalNominal(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		data["error"] = "Gagal menampilkan total pemasukan & pengeluaran " + err.Error()
	} else {
		data["total_pemasukan"] = pemasukan
		data["total_pengeluaran"] = pengeluaran
	}

	//TAMPILKAN DATA LIST
	listFinancial, err := models.NewFinancialModel(controller.db).FindAllFinancial(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		data["error"] = "Gagal menampilkan list keuangan " + err.Error()
	} else {
		data["financials"] = listFinancial
	}

	//FUNC IDR RUPIAH
	funcMap := template.FuncMap{
		"formatIDR": formatIDR,
		"indexNo":   func(a, b int) int { return a + b },
	}
	tmpl, _ := template.New(filepath.Base(layout)).Funcs(funcMap).ParseFiles(layout)
	tmpl.Execute(writer, data)
}
