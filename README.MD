#webserver
Simple own (without net/http) HTTP protocol implementation in Go

## How to use?
#### Example with Context usage: [\*Click\*](main.go)
```go
func main() {
	server := webserver.WebServer {
		Address: "localhost:1234", // Address for server
		WebsiteDir: "./example_website", // Working directory with website files
	}

	server.Route(test, "/test", "GET", "POST") // Route "/test" to function "root" for HTTP methods GET and POST
	server.Start()
}

func test(ctx *webserver.Context) {
    // Do something awesome
	ctx.WriteResponseFile(200, "test.html") // Return "test.html" to Client
}
```