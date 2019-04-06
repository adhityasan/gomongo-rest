package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adhityasan/gomongo-rest/config"
	"github.com/adhityasan/gomongo-rest/controller"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := mux.NewRouter()

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", config.Of.App.Host, config.Of.App.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Define Routes
	r.HandleFunc("/identify", controller.Identify).Methods("POST")
	r.HandleFunc("/identifyazure", controller.IdentifyByAzure).Methods("POST")
	r.HandleFunc("/go/aisatsu", controller.Aisatsu).Methods("GET")
	r.HandleFunc("/ocr", controller.DoOCR).Methods("POST")

	// Start Server
	go func() {

		log.Printf("Starting Server at %v...\n", srv.Addr)
		if errsrv := srv.ListenAndServe(); errsrv != nil {
			log.Fatal(errsrv)
		}

		errdb := dbConnection()
		if errdb != nil {
			log.Fatal(errdb)
		}

	}()

	// Graceful Shutdown
	waitShutdownSignal(srv)
}

func dbConnection() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, errCreateClient := mongo.NewClient(options.Client().ApplyURI(config.Of.Mongo.URL))

	if errCreateClient != nil {
		return errors.New("Fail create new mongo client")
	}

	errClientConnect := client.Connect(ctx)

	if errClientConnect != nil {
		return errors.New("Mongo client fail to connect")
	}

	return nil
}

func waitShutdownSignal(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down...")
	os.Exit(0)
}
