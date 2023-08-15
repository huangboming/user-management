package main

func main() {
	Server, _ := InitializeServer()
	Server.LoginMongo()
	Server.SetupRoute()
	Server.RunServer()
}
