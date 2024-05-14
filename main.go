package main

func main() {
	server := initWebServer()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
