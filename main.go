package main

import "log"

func main() {
	log.Println("Starting server...")
	server, err := NewServer()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server shutting down...")
}
