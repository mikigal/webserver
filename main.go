package main

import (
	"./webserver"
	"fmt"
)

func main() {
	server := webserver.WebServer{
		Address:    ":1234",
		WebsiteDir: "./example_website",
	}

	server.Route(test, "/test", "GET", "POST")
	server.Start()
}

func test(ctx *webserver.Context) {
	fmt.Println("===")
	fmt.Println("Method: " + ctx.Method)
	fmt.Println("Path: " + ctx.Path)
	fmt.Println("ClientAddress: " + ctx.ClientAddress.String())

	for key, value := range ctx.RequestHeaders {
		fmt.Println("RequestHeader: " + key + "=" + value)
	}

	for key, value := range ctx.UrlParams {
		fmt.Println("UrlParam: " + key + "=" + value)
	}

	for key, value := range ctx.PostParams {
		fmt.Println("PostParam: " + key + "=" + value)
	}
	fmt.Println("===")

	ctx.AddResponseHeader("Server", "ItsAwesomeWebServer")

	//ctx.WriteResponse(200, "You can return some text instead of file too!")
	//ctx.Redirect(http.StatusMovedPermanently, "/")
	ctx.WriteResponseFile(200, "test.html")
}
