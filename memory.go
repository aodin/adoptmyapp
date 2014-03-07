package main

import (
	"flag"
	"fmt"
	"github.com/kkochis/adoptmyapp/robot"
	"log"
	"net/http"
)

var port = flag.Int("port", 8081, "Server Port")

func main() {
	flag.Parse()

	_, err := robot.New(robot.NewMemoryDB())
	if err != nil {
		panic(err)
	}
	address := fmt.Sprintf(":%d", *port)
	log.Println("Running on address:", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
