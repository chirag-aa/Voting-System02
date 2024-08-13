package routes

import (
	"net/http"
	_ "voting-system/auth"
	"voting-system/vote"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/ws", vote.HandleConnection)

	return r
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Registration logic
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Login logic
}
