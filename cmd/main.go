package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"doe-base/idopost-backend/pkg/config"
	"doe-base/idopost-backend/pkg/controllers"

	"github.com/gorilla/mux"
)

func main() {

	//** routing variables
	route := mux.NewRouter()
	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		massage := []byte("welcome to Idopost")
		w.Write(massage)
	})
	route.HandleFunc("/mongodb-user", config.HandleMongodbConnectFormSubmit).Methods("POST", "OPTIONS")
	route.HandleFunc("/mongodb-user-url", config.HandleMongodbURLConnectFormSubmit).Methods("POST", "OPTIONS")
	route.HandleFunc("/mongodb-user-url-details", config.HandleMongodbURLConnectFormDetailsSubmit).Methods("POST", "OPTIONS")
	route.HandleFunc("/mongodb-get-all", controllers.GetAllRequest).Methods("GET", "OPTIONS")
	route.HandleFunc("/mongodb-get-connection-status", config.ConnectionStatus).Methods("GET", "OPTIONS")
	route.HandleFunc("/mongodb-log-out", config.HandleMongodbLogOut).Methods("PUT", "OPTIONS")

	route.HandleFunc("/api/post/handler", controllers.PostRequestHandle).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/delete-one/handler", controllers.DeleteOneRequestHandler).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/delete-list/handler", controllers.DeleteListRequestHandler).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/delete-all/handler", controllers.DeleteAllRequestHandler).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/get-one/handler", controllers.GetOneRequest).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/get-list/handler", controllers.GetListRequest).Methods("POST", "OPTIONS")

	os.Setenv("SEVER_PORT", ":8080")
	//** start the sever
	fmt.Printf("starting sever on port %v \n", os.Getenv("SEVER_PORT"))
	err := http.ListenAndServe(os.Getenv("SEVER_PORT"), route)
	if err != nil {
		log.Fatal("error listen: ", err)
	}
}
