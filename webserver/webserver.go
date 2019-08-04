package webserver

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (server *WebServer) Start() {
	if server.ErrorHandlers == nil {
		server.ErrorHandlers = make(map[int]func(ctx *Context))
	}

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
			if server.ErrorHandlers[http.StatusBadRequest] != nil {
				server.ErrorHandlers[http.StatusBadRequest](&ctx)
			} else {
				ctx.Error(http.StatusBadRequest, "Malformed request")
			}

			connection.Write(ctx.parseResponse())
			break
		}

		route, err := server.findRoute(ctx)
		if err == nil {
			route.Listener(&ctx)

			if ctx.ResponseCode == 0 {
				if server.ErrorHandlers[http.StatusInternalServerError] != nil {
					server.ErrorHandlers[http.StatusInternalServerError](&ctx)
				} else {
					ctx.Error(http.StatusInternalServerError, "Controller didn't write anything to context")
				}

				connection.Write(ctx.parseResponse())
				break
			}

			connection.Write(ctx.parseResponse())
		} else {
			exists := ctx.File(http.StatusOK, ctx.Path)
			if ctx.Path == "/" {
				exists = ctx.File(http.StatusOK, "index.html")
			}

			if exists == nil {
				connection.Write(ctx.parseResponse())
			} else {
				if server.ErrorHandlers[http.StatusNotFound] != nil {
					server.ErrorHandlers[http.StatusNotFound](&ctx)
				} else {
					ctx.Error(http.StatusNotFound, "")
				}

				connection.Write(ctx.parseResponse())
			}
		}
		break
	}
}

func (server *WebServer) Route(listener func(ctx *Context), path string, methods ...string) {
	server.Routes = append(server.Routes, Route{
		Path:     path,
		Methods:  methods,
		Listener: listener,
	})
}

func (server *WebServer) ErrorHandler(code int, listener func(ctx *Context)) {
	if server.ErrorHandlers == nil {
		server.ErrorHandlers = make(map[int]func(ctx *Context))
	}

	server.ErrorHandlers[code] = listener
}

func (ctx *Context) AddResponseHeader(name string, value string) {
	ctx.ResponseHeaders[name] = value
}

func (ctx *Context) Redirect(code int, target string) {
	ctx.ResponseCode = code
	ctx.AddResponseHeader("Location", target)
}

func (ctx *Context) CustomResponse(code int, mime string, content string) {
	ctx.ResponseType = mime
	ctx.ResponseCode = code
	ctx.ResponseBody = content
}

func (ctx *Context) File(code int, file string) error {
	bytes, err := ioutil.ReadFile(ctx.WebsiteDir + "/" + file)
	if err != nil {
		return err
	}

	ctx.CustomResponse(code, mime.TypeByExtension("."+strings.Split(file, ".")[len(strings.Split(file, "."))-1]), string(bytes))
	return nil
}

func (ctx *Context) JSON(code int, content string) {
	ctx.CustomResponse(code, "application/json; charset=utf8", content)
}

func (ctx *Context) HTML(code int, content string) {
	ctx.CustomResponse(code, "text/html; charset=utf-8", content)
}

func (ctx *Context) Error(code int, message string) {
	ctx.HTML(code, "<html><head><title>Error</title></head><body><h1>"+strconv.Itoa(code)+" "+http.StatusText(code)+"</h1><h2>"+message+"</h2></body></html>")
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
