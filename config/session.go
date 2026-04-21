package config

import (
	"time"

	"github.com/alexedwards/scs/v2"
)

var SessionManager *scs.SessionManager

func InitSession() {
	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour //SESSION AKAN BERAKHIR AFTER 24 JAM
	SessionManager.Cookie.Name = "financial_record_jan"
	SessionManager.Cookie.Path = "/"
	SessionManager.Cookie.HttpOnly = true
	SessionManager.Cookie.Secure = true
}
