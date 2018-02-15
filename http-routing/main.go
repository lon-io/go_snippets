package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

var validMethods = map[string]string{
	"get":  "GET",
	"post": "POST",
}

const (
	badRequestError = iota
	internalServerError
	pageNotFound
	resourceNotFound
)

func main() {
	port := 8080

	server, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer server.Close()

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		handleReq(conn)
	}
}

func handleReq(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	lineCounter := 0
	var method, path string

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line)

		if line == "" {
			break
		}

		parts := strings.Fields(line)

		if lineCounter == 0 {
			method = parts[0]
			path = parts[1]
		}

		lineCounter++
	}

	err := validateMethod(method)

	if err != nil {
		sendError(conn, err, badRequestError)
	} else {
		handleRequest(conn, method, path)
	}

}

func validateMethod(method string) error {
	fmt.Printf("Testing Method: %s\n", method)

	valid := false
	for _, v := range validMethods {
		if v == method {
			valid = true
		}
	}

	if valid == false {
		return errors.New("Method not supported")
	}

	return nil
}

func sendError(conn net.Conn, err error, errType int) {
	var status int
	var reason string

	response := err.Error()

	switch errType {
	case badRequestError:
		status = 400
		reason = "BadRequest"
	case resourceNotFound:
		status = 404
		reason = "ResourceNotFound"
	case pageNotFound:
		status = 404
		reason = "PageNotFound"
		response = `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>Page Not Found</strong></body></html>`
	case internalServerError:
	default:
		status = 500
		reason = "InternalServerError"
	}

	fmt.Printf("Sending error response: %v %v %v\n", status, reason, err.Error())

	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", status, reason)
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(response))
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, response)
}

func handleRequest(conn net.Conn, method, path string) {
	if method == validMethods["get"] {
		switch path {
		case "/":
			sendSuccess(conn, `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>Home</strong></body></html>`)
		case "/about":
			sendSuccess(conn, `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>About</strong></body></html>`)
		case "/contact":
			sendSuccess(conn, `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>Contact Us</strong></body></html>`)
		default:
			sendError(conn, errors.New("Page Not Found"), pageNotFound)
		}
	} else {
		switch path {
		case "/user":
			sendSuccess(conn, "User created")
		case "/post":
			sendSuccess(conn, "Post created")
		default:
			sendError(conn, fmt.Errorf("Route %s is not registered", path), resourceNotFound)
		}
	}
}

func sendSuccess(conn net.Conn, body string) {
	fmt.Printf("Sending success response: %v\n", body)
	fmt.Fprintf(conn, "HTTP/1.1 200 Ok\r\n")
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, body)
}
