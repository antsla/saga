package transport

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() Server {
	s := Server{}
	s.router = mux.NewRouter()

	return s
}

func (s Server) Start() error {
	return http.ListenAndServe(":"+os.Getenv("HTTP_BIND"), s.router)
}
