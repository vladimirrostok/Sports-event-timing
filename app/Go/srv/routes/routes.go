package routes

import (
	checkpoint_controller "sports/backend/srv/controllers/checkpoint"
	result_controller "sports/backend/srv/controllers/result"
	sportsmen_controller "sports/backend/srv/controllers/sportsmen"
	"sports/backend/srv/middleware"
	"sports/backend/srv/server"
)

func InitializeRoutes(s *server.Server) {
	s.Router.HandleFunc("/dashboard", s.Dashboard.ResultsHandler)

	s.Router.HandleFunc("/results", middleware.SetMiddlewareJSON(result_controller.AddResult(s))).Methods("POST")
	s.Router.HandleFunc("/results", middleware.SetMiddlewareJSON(result_controller.GetLastTenResults(s))).Methods("GET")
	s.Router.HandleFunc("/finish", middleware.SetMiddlewareJSON(result_controller.AddFinishTime(s))).Methods("POST")
	s.Router.HandleFunc("/checkpoints", middleware.SetMiddlewareJSON(checkpoint_controller.AddCheckpoint(s))).Methods("POST")
	s.Router.HandleFunc("/sportsmens", middleware.SetMiddlewareJSON(sportsmen_controller.AddSportsmen(s))).Methods("POST")
}
