package main

import (
	"financial_record/config"
	"financial_record/routes"

	"log"
	"net/http"
)

func main() {
	//http.HandleFunc("/testing", func(writer http.ResponseWriter, request *http.Request) {
	//	tmpl, err := template.ParseFiles("views/testing.html")
	//	if err != nil {
	//		http.Error(writer, "Template error", http.StatusInternalServerError)
	//		return
	//	}
	//	//KIRIM DATA DARI GO KE HTML
	//	data := make(map[string]interface{})
	//	data["name"] = "Fajar Setyo Pambudi"
	//	data["training"] = "Golang"
	//	data["duration"] = 4
	//
	//	if request.Method == http.MethodPost {
	//		//AMBIL DATA DARI HTML
	//		request.ParseForm()
	//		username := request.FormValue("username")
	//		age := request.FormValue("age")
	//		fmt.Println(username, age)
	//	}
	//	tmpl.Execute(writer, data)
	//
	//})
	//http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
	//	tmpl, err := template.ParseFiles("views/auth/login.html")
	//	if err != nil {
	//		http.Error(writer, "Template error", http.StatusInternalServerError)
	//		return
	//	}
	//	if request.Method == http.MethodPost {
	//		//AMBIL DATA DARI HTML
	//		request.ParseForm()
	//		username := request.FormValue("username")
	//		age := request.FormValue("age")
	//		fmt.Println(username, age)
	//	}
	//	tmpl.Execute(writer, nil)
	//
	//})

	config.InitSession()
	mux := http.NewServeMux()

	//DAFTARKAN DIRECTORI
	mux.Handle("/user_photo/", http.StripPrefix("/user_photo/", http.FileServer(http.Dir("public/user_photo"))))

	//PANGGIL DATABASE
	db := config.InitDatabase()

	//PANGGIL ROUTE
	routes.Route(db, mux)

	log.Println("Service sedang berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", config.SessionManager.LoadAndSave(mux))
}
