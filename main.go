package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/abdelrhman-basyoni/goresp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: redis-cli-go <host> <port>")
		return
	}

	host := os.Args[1]
	port := os.Args[2]
	address := net.JoinHostPort(host, port)

	// Connect to the Redis server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Failed to connect to Redis server at %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to Redis server at %s\n", address)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println()
		fmt.Print("redis> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Exit command
		if strings.EqualFold(line, "exit") {
			fmt.Println("Exiting...")
			break
		}

		// Send the command to Redis server
		_, err := conn.Write(goresp.SerializeCommand(line))
		if err != nil {
			fmt.Printf("Error writing to Redis server: %v\n", err)
			break
		}

		// Read the response from Redis server
		resReader := goresp.NewRespReader(conn)
		res, err := resReader.Read()
		if err != nil {
			fmt.Printf("Error reading  from server: %v\n", err)
			break
		}
		fmt.Println(goresp.SerializeValue(res))
	}
}
