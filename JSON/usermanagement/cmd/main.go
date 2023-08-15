package main

func main() {
	Server, _ := InitializeServer()
	Server.SetupRoute()
	Server.RunServer()
}
