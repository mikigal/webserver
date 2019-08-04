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
	server.ErrorHandler(500, customError)
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

	//ctx.Redirect(http.StatusMovedPermanently, "/anotherPath")
	//ctx.HTML(200, "<html><head></head><body><h1>Test</h1></body></html>")
	//ctx.JSON(200, "{\"test\": \"abc\"}")
	//ctx.Error(403, "")
	ctx.File(200, "test.html")
}

func customError(ctx *webserver.Context) {
	ctx.JSON(404, "{\"error\": \"not_found\"}")
}
