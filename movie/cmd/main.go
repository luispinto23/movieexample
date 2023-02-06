package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/luispinto23/movieexample/movie/internal/controller/movie"
	metadatagateway "github.com/luispinto23/movieexample/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/luispinto23/movieexample/movie/internal/gateway/rating/http"
	httphandler "github.com/luispinto23/movieexample/movie/internal/handler/http"
	"github.com/luispinto23/movieexample/pkg/discovery"
	"github.com/luispinto23/movieexample/pkg/discovery/consul"
)

var serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Starting the %s service on port %d", serviceName, port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	// span a go func to report the health state every second
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle(fmt.Sprintf("/%s", serviceName), http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
