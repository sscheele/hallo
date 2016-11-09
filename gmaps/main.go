package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey       string
	origin       = "Perth"         //place of origin
	destination  = "Sydney"        //destination coordinates/address
	arrivalTime  = "1480543258"    //arrival time (seconds since epoch)
	alternatives = false           //find alternative routes
	avoid        = "tolls|ferries" //valid values are tolls, highways, ferries, | separated
	trafficModel = "best_guess"    //valid values are optimistic, best_guess, and pessimistic.
)

func check(err error, description string) {
	if err != nil {
		fmt.Printf("fatal error in %s: %s", description, err)
	}
}

func main() {
	f, err := os.Open("api-key.txt")
	check(err, "opening file")
	reader := bufio.NewReader(f)
	apiKey, _ = reader.ReadString('\n')
	apiKey = apiKey[:len(apiKey)-1] //remove trailing newline (required because it's a text file)
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(2))
	check(err, "new maps client")

	r := &maps.DirectionsRequest{
		Origin:       origin,
		Destination:  destination,
		ArrivalTime:  arrivalTime,
		Alternatives: alternatives,
		Mode:         maps.TravelModeDriving,
	}

	lookupTrafficModel(trafficModel, r)

	if avoid != "" {
		lookupAvoidPoints(avoid, r)
	}

	routes, waypoints, err := client.Directions(context.Background(), r)
	check(err, "getting directions")

	fmt.Println(waypoints)
	fmt.Println(routes)
}

func lookupTrafficModel(trafficModel string, r *maps.DirectionsRequest) {
	switch trafficModel {
	case "optimistic":
		r.TrafficModel = maps.TrafficModelOptimistic
	case "best_guess":
		r.TrafficModel = maps.TrafficModelBestGuess
	case "pessimistic":
		r.TrafficModel = maps.TrafficModelPessimistic
	case "":
		// ignore
	default:
		log.Fatalf("Unknown traffic mode %s", trafficModel)
	}
}

func lookupAvoidPoints(avoidPts string, r *maps.DirectionsRequest) {
	for _, a := range strings.Split(avoidPts, "|") {
		switch a {
		case "tolls":
			r.Avoid = append(r.Avoid, maps.AvoidTolls)
		case "highways":
			r.Avoid = append(r.Avoid, maps.AvoidHighways)
		case "ferries":
			r.Avoid = append(r.Avoid, maps.AvoidFerries)
		}
	}
}
