package webserver

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"strings"
	"time"
)

func (server *WebServer) Start() {
	listener, _ := net.Listen("tcp4", server.Address)
	for {
		connection, _ := listener.Accept()
		go server.handle(connection)
	}
}

func (server *WebServer) handle(connection net.Conn) {
	connection.SetReadDeadline(time.Now().Add(2 * time.Minute))
	buffer := make([]byte, 4096)
	defer connection.Close()

	for {
		connection.Read(buffer)

		ctx, err := server.parseContext(connection, string(buffer))
		if err != nil {
			errorCtx := Context{
				ResponseType: "text/html; charset=utf8",
				ResponseCode: http.StatusBadRequest,
				ResponseBody: "<html><head><title>Error</title></head><body><h1>400 Bad Request</h1></body></html>",
			}

			connection.Write(errorCtx.parseResponse())
			break
		}

		route, err := server.findRoute(ctx)
		if err == nil {
			route.Listener(&ctx)
			connection.Write(ctx.parseResponse())
		} else {
			exists := ctx.WriteResponseFile(http.StatusOK, ctx.Path)
			if ctx.Path == "/" {
				exists = ctx.WriteResponseFile(http.StatusOK, "index.html")
			}

			if exists == nil {
				connection.Write(ctx.parseResponse())
			} else {
				ctx.WriteResponse(http.StatusNotFound, "<html><head><title>Error</title></head><body><h1>404 Not Found</h1></body></html>")
				connection.Write(ctx.parseResponse())
			}
		}

		break
	}
}

func (server *WebServer) Route(listener func(request *Context), path string, methods ...string) {
	server.Routes = append(server.Routes, Route{
		Path:     path,
		Methods:  methods,
		Listener: listener,
	})
}

func (server *WebServer) findRoute(ctx Context) (Route, error) {
	for _, route := range server.Routes {
		if route.Path == ctx.Path {
			for _, method := range route.Methods {
				if ctx.Method == method {
					return route, nil
				}
			}
		}
	}

	return Route{}, errors.New("Route not found")
}

func (ctx *Context) WriteResponseFile(code int, file string) error {
	bytes, err := ioutil.ReadFile(ctx.WebsiteDir + "/" + file)
	if err != nil {
		return err
	}

	ctx.ResponseCode = code
	ctx.ResponseBody = string(bytes)
	ctx.ResponseType = mime.TypeByExtension("." + strings.Split(file, ".")[len(strings.Split(file, "."))-1])
	return nil
}

func (ctx *Context) WriteResponse(code int, content string) {
	ctx.ResponseCode = code
	ctx.ResponseBody = content
	ctx.ResponseType = "text/html; charset=utf-8"
}

func (ctx *Context) AddResponseHeader(name string, value string) {
	ctx.ResponseHeaders[name] = value
}

type WebServer struct {
	Address    string
	Routes     []Route
	WebsiteDir string
}

type Context struct {
	Method          string
	Path            string
	ClientAddress   net.Addr
	UrlParams       map[string]string
	PostParams      map[string]string
	WebsiteDir      string
	ResponseCode    int
	ResponseBody    string
	ResponseType    string
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
}

type Route struct {
	Path     string
	Methods  []string
	Listener func(ctx *Context)
}
