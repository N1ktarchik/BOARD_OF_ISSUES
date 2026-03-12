package tarnsport

import (
	"Board_of_issuses/internal/features/transport/http/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Handlers *handlers.UserHandler
	Auth     *AuthHandler
}

func NewHTTPServer(Handlers *handlers.UserHandler, Auth *AuthHandler) *HTTPServer {
	return &HTTPServer{
		Handlers: Handlers,
		Auth:     Auth,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	//главная страница
	router.Path("/").HandlerFunc(s.Handlers.HandleBase)

	router.Path("/register").Methods("POST").HandlerFunc(s.Handlers.HandleCreateUser)
	router.Path("/login").Methods("POST").HandlerFunc(s.Handlers.HandleLoginUser)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(s.Auth.AuthMiddleware)

	api.Path("/users/name").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeUserName)
	api.Path("/users/password").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeUserPassword)
	api.Path("/users/email").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeUserEmail)
	api.Path("/users").Methods("POST").HandlerFunc(s.Handlers.HandleConnectUserToDesk)

	api.Path("/desks").Methods("POST").HandlerFunc(s.Handlers.HandleCreateDesk)
	api.Path("/desks/{id}/name").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeDeskName)
	api.Path("/desks/{id}/password").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeDeskPassword)
	api.Path("/desks/{id}/owner").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeDeskOwner)
	api.Path("/desks/{id}").Methods("DELETE").HandlerFunc(s.Handlers.HandleDeleteDesk)
	api.Path("/desks").Methods("GET").HandlerFunc(s.Handlers.HandleGetAllDesksId)

	api.Path("/tasks").Methods("POST").HandlerFunc(s.Handlers.HandleCreateTask)
	api.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(s.Handlers.HandleDeleteTask)
	api.Path("/tasks/{id}/complyte").Methods("PATCH").HandlerFunc(s.Handlers.HandleComplyteTask)
	api.Path("/tasks/{id}/time").Methods("PATCH").HandlerFunc(s.Handlers.HandleAddTimeToTask)
	api.Path("/tasks/{id}/description").Methods("PATCH").HandlerFunc(s.Handlers.HandleChangeTaskDescription)
	api.Path("/tasks/{deskId}").Methods("GET").HandlerFunc(s.Handlers.HandleGetAllTasks)
	api.Path("/tasks").Methods("GET").Queries("done", "{done}", "desk_id", "{desk_id}").HandlerFunc(s.Handlers.HandleGetTasksWithParams)

	return http.ListenAndServe(":8080", router)
}
