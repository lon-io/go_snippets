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
		handleError(conn, err, badRequestError)
	} else {
		handleSuccess(conn, method, path)
	}

}

func validateMethod(method string) error {
	valid := false
	for v := range validMethods {
		if v == method {
			valid = true
		}
	}

	if valid == false {
		return errors.New("Method not supported")
	}

	return nil
}

func handleError(conn net.Conn, err error, errType int) {
	var status int
	var reason string

	switch errType {
	case badRequestError:
		status = 400
		reason = "BadRequest"
	case internalServerError:
		status = 500
		reason = "InternalServerError"
	}

	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", status, reason)
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, "")
}

func handleSuccess(conn net.Conn, method, path string) {
	if method == validMethods["get"] {
		switch path {
		case "/":
			sendSuccess(conn, `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>Home</strong></body></html>`)
		case "/about":
			sendSuccess(conn, `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body><strong>About</strong></body></html>`)
		}
	} else {
		switch path {
		case "/user":
			sendSuccess(conn, "User created")
		case "/post":
			sendSuccess(conn, "Post created")
		}
	}
}

func sendSuccess(conn net.Conn, body string) {
	fmt.Fprintf(conn, "HTTP/1.1 200 Ok\r\n")
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, body)
}
