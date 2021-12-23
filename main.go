package main

import (
	api "github.com/asad1123/url-shortener/server/routes"
)

func main() {
	var server api.Routes
	server.startGin()
}
