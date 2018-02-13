package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	port := 8080

	server, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer server.Close()

	for {
		requestConn, err := server.Accept()

		if err != nil {
			log.Fatalln(err)
			continue
		}

		handleRequest(requestConn)
	}
}

func handleRequest(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	var host, path string

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line)

		// Check for ending empty line
		if line == "" {
			break
		}

		parts := strings.Fields(line)

		switch parts[0] {
		case "GET":
			path = parts[1]
		case "Host:":
			host = parts[1]
		}
	}

	uri := fmt.Sprintf("%v%v", host, path)

	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(uri))
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, uri)
}
