package routes

import "sports/backend/srv/server"

func InitializeRoutes(s *server.Server) {
	s.Router.HandleFunc("/dashboard", s.Dashboard.ResultsHandler)
}
