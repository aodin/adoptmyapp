package main

import (
	"flag"
	"fmt"
	"github.com/kkochis/adoptmyapp/robot"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

var port = flag.Int("port", 8081, "Server Port")

func main() {
	flag.Parse()

	// Start the mongodb session
	session, err := mgo.Dial(`localhost:27017`)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	_, err = robot.New(robot.NewMongoDB(session, "test"))
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
