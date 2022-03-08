package route

import (
	"log"
	"net/http"
)

func HandleClientRoutes(res http.ResponseWriter, req *http.Request) {
	log.Print(req.URL.Path)
	if req.URL.Path == "/" {
		// middleware.Login()
	}
}
