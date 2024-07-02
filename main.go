package main

import (
	"smoothstart/server"
)

func main() {
	s := server.NewSSS()
	s.Server.Logger.Fatal(s.Server.Start(":8080"))
}
