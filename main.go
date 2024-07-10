package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/login_example")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registerHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !isSessionValid(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Welcome to the secure area!")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		user := r.FormValue("username")
		pass := r.FormValue("password")
		if isValidUser(user, pass) {
			sessionID := createSession(user)
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: sessionID,
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		fmt.Fprintf(w, "Invalid credentials")
		return
	}

	fmt.Fprintf(w, `<html><body><form method="post" action="/login">
                    Username: <input type="text" name="username"><br>
                    Password: <input type="password" name="password"><br>
                    <input type="submit" value="Login">
                    </form></html></body>`)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		deleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		user := r.FormValue("username")
		pass := r.FormValue("password")
		



		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}
		fmt.Printf(string(hashedPassword))

		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user, hashedPassword)
		
		if err != nil {
			http.Error(w, "Username is already taken.", 400)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	fmt.Fprintf(w, `<html><body><form method="post" action="/register">
                    Username: <input type="text" name="username"><br>
                    Password: <input type="password" name="password"><br>
                    <input type="submit" value="Register">
                    </form></html></body>`)
}

func isValidUser(username, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

var sessions = map[string]string{}

func createSession(username string) string {
	sessionID := fmt.Sprintf("%d", len(sessions)+1)
	sessions[sessionID] = username
	return sessionID
}

func isSessionValid(sessionID string) bool {
	_, exists := sessions[sessionID]
	return exists
}

func deleteSession(sessionID string) {
	delete(sessions, sessionID)
}
