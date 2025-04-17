package auth

import (
	"database/sql"
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
	"real_time_forum/backend/models"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage renders the login page, or redirects if the user is already authenticated.
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		controllers.RenderError(w, http.StatusBadRequest)
		return
	}

	Email := r.FormValue("Email")
	Pass := r.FormValue("Password")

	ID, _, err := models.VerifyEmail(db, Email)
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}

	if ID == -1 {
		// e := models.ErrorRegister{ErrEmail: "Incorrect email"}
		// controllers.RenderTemplate(w, "login.html", e, http.StatusConflict)
		// return
	}

	PasswordDatabase, err := models.GetPassword(db, int(ID))
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(PasswordDatabase), []byte(Pass))
	if err != nil {
		// e := models.ErrorRegister{ErrPassword: "Incorrect Password"}
		// controllers.RenderTemplate(w, "login.html", e, http.StatusConflict)
		// return
	}

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		controllers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := middleware.VerifyCookie(r, db)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = controllers.RenderTemplate(w, "login.html", nil, http.StatusOK)
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}
}
