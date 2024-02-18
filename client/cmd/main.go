package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	serverAddr    = "wowsrv:8080"
	challengeSize = 64
	difficulty    = 16 // Must match the server's difficulty setting
)

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server. Receiving challenge...")

	challenge := make([]byte, challengeSize)
	_, err = bufio.NewReader(conn).Read(challenge)
	if err != nil {
		fmt.Println("Error reading challenge:", err.Error())
		return
	}

	fmt.Println("Solving challenge...")
	solution := solveChallenge(challenge)

	fmt.Println("Sending solution...")
	_, err = conn.Write(solution)
	if err != nil {
		fmt.Println("Error sending solution:", err.Error())
		return
	}

	response, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			// EOF after a write is expected; this means the server has sent the quote and closed the connection.
			fmt.Println("Received quote:", string(response))
		} else {
			fmt.Println("Error reading server response:", err.Error())
		}
		return
	}
}

func solveChallenge(challenge []byte) []byte {
	var solution []byte
	var hash [32]byte

	for {
		solution = make([]byte, challengeSize)
		_, err := rand.Read(solution)
		if err != nil {
			fmt.Println("Error generating solution:", err)
			continue
		}

		combined := append(challenge, solution...)
		hash = sha256.Sum256(combined)

		if isValidSolution(hash[:]) {
			break
		}
	}

	return solution
}

func isValidSolution(hash []byte) bool {
	hashString := hex.EncodeToString(hash)
	for i := 0; i < difficulty/4; i++ {
		if hashString[i] != '0' {
			return false
		}
	}
	return true
}
