package webserver

import (
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (server *WebServer) parseContext(connection net.Conn, raw string) (Context, error) {
	ctx := Context{
		UrlParams:       make(map[string]string),
		PostParams:      make(map[string]string),
		ResponseHeaders: make(map[string]string),
		RequestHeaders:  make(map[string]string),
		ClientAddress:   connection.RemoteAddr(),
		WebsiteDir:      server.WebsiteDir,
	}

	for index, line := range strings.Split(raw, "\r\n") {
		if index != 0 && strings.Contains(line, ": ") {
			header := strings.Split(strings.Replace(line, "\r\n", "", 1), ": ")
			if len(header) < 2 {
				return Context{}, errors.New("400")
			}

			ctx.RequestHeaders[header[0]] = header[1]
		}
	}

	spaces := strings.Split(raw, " ")
	ctx.Method = spaces[0]

	path := spaces[1]
	ctx.Path = strings.Split(path, "?")[0]

	if !strings.Contains(path, "?") && strings.Contains(path, "&") {
		return Context{}, errors.New("400")
	}

	if strings.Contains(path, "?") {
		rawParam := strings.Split(strings.Split(path, "?")[1], "&")[0]
		if !strings.Contains(rawParam, "=") {
			return Context{}, errors.New("400")
		}

		parseParam(ctx.UrlParams, rawParam)
	}

	for index, rawParam := range strings.Split(path, "&") {
		if index != 0 { // Ignore real path and "?" param
			if !strings.Contains(rawParam, "=") {
				return Context{}, errors.New("400")
			}

			parseParam(ctx.UrlParams, rawParam)
		}
	}

	body := strings.Split(raw, "\r\n\r\n")
	if ctx.Method == "POST" && len(body) == 2 {
		for _, rawParam := range strings.Split(body[1], "&") {
			if !strings.Contains(rawParam, "=") {
				return Context{}, errors.New("400")
			}

			parseParam(ctx.PostParams, rawParam)
		}
	}

	return ctx, nil
}

func (ctx *Context) parseResponse() []byte {
	var response = "HTTP/1.1 " + strconv.Itoa(ctx.ResponseCode) + " " + http.StatusText(ctx.ResponseCode) + "\r\n "

	for name, value := range ctx.ResponseHeaders {
		parseHeader(&response, name, value)
	}

	parseHeader(&response, "Content-Type", ctx.ResponseType)
	parseHeader(&response, "Date", time.Now().Format(time.RFC1123))
	parseHeader(&response, "Server", "go-webserver")
	parseHeader(&response, "Connection", "Keep-Alive")
	parseHeader(&response, "Keep-Alive", "timeout=15, max=100")
	response += "\r\n" + ctx.ResponseBody

	return []byte(response)
}

func parseHeader(response *string, name string, value string) {
	if !strings.Contains(*response, name+": ") {
		*response += name + ": " + value + "\r\n"
	}
}

func parseParam(params map[string]string, raw string) {
	split := strings.Split(raw, "=")
	params[split[0]] = split[1]
}
