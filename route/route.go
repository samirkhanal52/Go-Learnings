package route

import (
	"log"
	"net/http"

	"github.com/samirkhanal52/go-cli-chat-app/middleware"
)

func HandleClientRoutes(res http.ResponseWriter, req *http.Request) {
	log.Print(req.URL.Path)
	if req.URL.Path == "/" {
		middleware.Login()
	}
}
