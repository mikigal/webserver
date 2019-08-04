# webserver
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
    // If there isn't route for URL server serve file from WebsiteDir, for / it's index.html
    server.Start() // Start server!
}

func test(ctx *webserver.Context) {
    // Do something awesome

    //ctx.Redirect(http.StatusMovedPermanently, "/anotherPath")
    //ctx.HTML(200, "<html><head></head><body><h1>Test</h1></body></html>")
    //ctx.JSON(200, "{\"test\": \"abc\"}")
    //ctx.Error(403, "")
    ctx.File(200, "test.html")
}
```