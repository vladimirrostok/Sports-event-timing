package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"sports/backend/srv/controllers/dashboard"
)

// Server is a wrapper for the service context.
type Server struct {
	Dashboard *dashboard_controller.Dashboard
	DB        *gorm.DB
	Router    *mux.Router
	Addr      string
}
