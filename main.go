package main

import "github.com/go-jarvis/herts/server"

func main() {

	s := &server.Config{
		Listen: ":8081",
	}

	err := s.Run()
	if err != nil {
		panic(err)
	}
}
