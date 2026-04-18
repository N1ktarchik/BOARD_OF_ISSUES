package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"Board_of_issuses/docs"
	
)

type Server struct {
	Router *mux.Router
	log    *slog.Logger
}

func NewServer(log *slog.Logger) *Server {
	return &Server{
		Router: mux.NewRouter(),
		log:    log,
	}
}

func (s *Server) Run(addr string) error {
	s.log.Info("starting server", slog.String("addr", addr))

	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		s.log.Error("server failed to start or crashed",
			slog.String("addr", addr),
			slog.Any("err", err),
		)
		return err
	}

	return nil
}

func (s *Server) RegisterSwagger() {

	s.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)).Methods("GET")

	s.Router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	}).Methods("GET")

	s.log.Info("swagger have registred")
}
