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
		_, err := conn.Write(serializeCommand(line))
		if err != nil {
			fmt.Printf("Error writing to Redis server: %v\n", err)
			break
		}

		// Read the response from Redis server
		// response, err := bufio.NewReader(conn).ReadString('\n')
		resReader := goresp.NewRespReader(conn)
		// if err != nil {
		// 	fmt.Printf("Error reading response from Redis server: %v\n", err)
		// 	break
		// }

		res, err := resReader.Read()
		if err != nil {
			fmt.Printf("Error reading  from server: %v\n", err)
			break
		}
		fmt.Println("response from Redis server")
		fmt.Println(serializeReaderCommand(res))
	}
}

// serializeCommand converts the command into RESP (REdis Serialization Protocol) format
func serializeCommand(command string) []byte {
	parts := strings.Fields(command)

	var test []goresp.Value

	for _, part := range parts {

		bulk := goresp.Value{Typ: "bulk", Bulk: part}
		test = append(test, bulk)
	}

	newCommand := goresp.Value{Typ: "array", Array: test}

	fmt.Println(len(newCommand.Array))
	fmt.Println(newCommand)
	fmt.Println(newCommand.Marshal())
	return newCommand.Marshal()
}

func serializeReaderCommand(res goresp.Value) string {

	switch res.Typ {
	case "bulk":
		{
			return res.Bulk
		}

	default:
		{
			return res.Bulk
		}
	}
}
