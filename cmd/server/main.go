package main

import (
	"context"
	"flag"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/db"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/handler"
)

func main() {
	var dburl = flag.String("db", "postgres://postgres:infoblox123@172.17.0.2:5432/postgres?sslmode=disable", "database connection url")

	flag.Parse()

	mux := http.NewServeMux()

	con := db.InitDb(*dburl)
	resSvc := service.NewResource(db.NewResource(con))

	resHandler := handler.NewResource(resSvc)
	clusterHandler := handler.NewCluster(resSvc)

	mux.Handle(handler.AddResourceURL, http.HandlerFunc(resHandler.HandleAddRequest))
	mux.Handle(handler.ClusterGetPointsURL, http.HandlerFunc(clusterHandler.HandleGetPointsRequest))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("server starting on localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			log.Println("server error: ", err.Error())
			shutdown(server)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	shutdown(server)

}

func shutdown(server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("server shutdown failed:", err)
	}
	log.Printf("server exited properly")
}
