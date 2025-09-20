package main

import (
	"fmt"
	"log"
	"net"

	request "tcp.to.http/internal/requests"
)

func main() {
	// f, err := os.Open("message.txt")
	// if err != nil {
	// 	log.Fatal("crash", "crash", err)
	// }

	listener, err := net.Listen("tcp", ":42068")

	if err != nil {
		log.Fatal("Error", "Error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error", "Error", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("Error", "Error", err)
		}

		fmt.Printf("Request line: \n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Printf("Headers: \n")
		r.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
		fmt.Printf("Body: \n")
		fmt.Printf("%s \n", r.Body)
	}

}
