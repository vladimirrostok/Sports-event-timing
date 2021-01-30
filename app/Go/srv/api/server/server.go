//mockgen -destination=mocks/http_client.go -package=mocks go-ddd-cqrs-example/usersapi/context HTTPClient
package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
)

// HTTPClient interface to mock the network requests for test purposes.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// Server is a wrapper for the service context.
type Server struct {
	DB         *gorm.DB
	Router     *mux.Router
	HTTPClient HTTPClient
	Addr       string
}
