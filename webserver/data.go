package webserver

import "net"

type WebServer struct {
	Address       string
	Routes        []Route
	WebsiteDir    string
	ErrorHandlers map[int]func(ctx *Context)
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
