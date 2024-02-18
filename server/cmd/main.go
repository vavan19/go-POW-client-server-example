package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
)

const (
	listenAddr    = ":8080"
	difficulty    = 16 // Number of bits that must be zero (the higher, the more difficult)
	challengeSize = 64
)

var quotes = []string{
	"The only way to do great work is to love what you do. - Steve Jobs",
	"The best way out is always through",
	"Always Do What You Are Afraid To Do",
	"The journey of a thousand miles begins with one step",
}

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on " + listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	challenge := make([]byte, challengeSize)
	_, err := rand.Read(challenge)
	if err != nil {
		fmt.Println("Error generating challenge:", err)
		return
	}

	conn.Write(challenge) // Send the challenge to the client

	solution := make([]byte, challengeSize)
	_, err = bufio.NewReader(conn).Read(solution)
	if err != nil {
		fmt.Println("Error reading solution:", err)
		return
	}

	if validateSolution(challenge, solution) {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(quotes))))
		conn.Write([]byte(quotes[randomIndex.Int64()]))
	} else {
		conn.Write([]byte("Invalid PoW"))
	}
}

func validateSolution(challenge, solution []byte) bool {
	combined := append(challenge, solution...)
	hash := sha256.Sum256(combined)
	hashString := hex.EncodeToString(hash[:])

	// Check if the hash has the required number of leading zeros
	for i := 0; i < difficulty/4; i++ {
		if hashString[i] != '0' {
			return false
		}
	}
	return true
}
