package routes

import (
	"database/sql"
	"financial_record/controllers"
	"financial_record/helpers"
	"net/http"
)

func Route(db *sql.DB, mux *http.ServeMux) {
	//AUTH CONTROLLER
	authController := controllers.NewAuthController(db)
	mux.HandleFunc("/register", helpers.GuestOnly(authController.Register))
	mux.HandleFunc("/login", helpers.GuestOnly(authController.Login))
	mux.HandleFunc("/logout", helpers.AuthOnly(authController.Logout))

	//FINANCIAL CONTROLLER
	financialController := controllers.NewFinancialController(db)
	mux.HandleFunc("/home", helpers.AuthOnly(financialController.Home))
	mux.HandleFunc("/financial/add_financial_record", helpers.AuthOnly(financialController.AddFinancialRecord))
	mux.HandleFunc("/financial/edit_financial_record", helpers.AuthOnly(financialController.EditFinancialRecord))
	mux.HandleFunc("/financial/delete_financial_record", helpers.AuthOnly(financialController.DeleteFinancialRecord))
	mux.HandleFunc("/financial/download_financial_record", helpers.AuthOnly(financialController.DownloadFinancialRecord))

	//PROFILE
	userController := controllers.NewUserController(db)
	mux.HandleFunc("/profile", helpers.AuthOnly(userController.Profile))

}
