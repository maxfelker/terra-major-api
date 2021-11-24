package server

import (
	"log"
	"net/http"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type Routes []Route

func Start(PORT string) {
	logMessage := "Listen for requests at http://localhost" + PORT
	log.Println(logMessage)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func RegisterRoutes(routes []Route) {
	log.Println("Registering routes...")
	for _, route := range routes {
		http.HandleFunc(route.Path, route.Handler)
		log.Println(route.Method + " " + route.Path)
	}
}

/*
func validateRequest(request *http.Request) bool {
	for _, route := range routes {
		if request.URL.Path == route.path {
			if request.Method == route.method {
				return true
			}
		}
	}

	return false
}

func requestFailed(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.Method + " " + request.URL.Path + " is an invalid endpoint")
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	writer.Write(nil)
}

*/
