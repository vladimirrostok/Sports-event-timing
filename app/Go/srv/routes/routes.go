package routes

import (
	result_controller "sports/backend/srv/controllers/result"
	"sports/backend/srv/middleware"
	"sports/backend/srv/server"
)

func InitializeRoutes(s *server.Server) {
	s.Router.HandleFunc("/dashboard", s.Dashboard.ResultsHandler)
	s.Router.HandleFunc("/result", middleware.SetMiddlewareJSON(result_controller.AddResult(s))).Methods("POST")
}
