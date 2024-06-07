package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Pokemon struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Types []string          `json:"types"`
	Stats map[string]string `json:"stats"`
	Exp   string            `json:"exp"`
}

type Player struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	PokeBalls []Pokemon `json:"pokeBalls"`
}

var POKEMONS []Pokemon

// var exchangeMap = make(map[string]string)
// var exchangeJSON []byte
var PLAYERS []Player
var ROWS, COLS = 1000, 1000
var BOARD = make([][]string, ROWS)
var CONNECTIONS = make(map[string]net.Conn)

// ////////////////////////////////////////////////////////////////////////////////////
// Load players.json to PLAYERS
func loadPlayers(filename string) ([]Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var players []Player
	err = json.Unmarshal(bytes, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func loadPokemons(filename string) ([]Pokemon, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var pokemons []Pokemon
	err = json.Unmarshal(bytes, &pokemons)
	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

//////////////////////////////////////////////////////////////////////////////////////

func verifyPlayer(username, password string, players []Player) bool {
	for _, user := range players {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			return err == nil
		}
	}
	return false
}

// HandleConnection handles incoming client connections
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a reader to read data from the connection
	reader := bufio.NewReader(conn)
	for {
		// Read data from the connection
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		// Print the message received from the client
		fmt.Printf("Message received: %s", message)

		// Echo the message back to the client

		for _, tcpConn := range CONNECTIONS {
			if conn != tcpConn {
				tcpConn.Write([]byte("Echo: " + message))
			}
		}
	}
}

func main() {
	// Set up BOARD
	for i := range BOARD {
		BOARD[i] = make([]string, COLS)
	}
	// Load pokemon.json to array POKEMONS
	POKEMONS, err := loadPokemons("pokedex.json")
	checkError(err)

	//Load players.json to array PLAYERS
	PLAYERS, err := loadPlayers("players.json")
	checkError(err)

	fmt.Println(PLAYERS[0].Username)
	for range 50 {
		newPokemonId := rand.Intn(len(POKEMONS)) + 1
		spawnX := rand.Intn(ROWS)
		spawnY := rand.Intn(COLS)
		BOARD[spawnX][spawnY] = strconv.Itoa(newPokemonId)
	}

	// Start listening for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Verfify username and password
		infoReader := bufio.NewReader(conn)

		// Get username
		username, err := infoReader.ReadString('\n')
		checkError(err)

		// Get password
		password, err := infoReader.ReadString('\n')
		checkError(err)
		fmt.Println(strings.TrimSpace(username), strings.TrimSpace(password))
		if verifyPlayer(strings.TrimSpace(username), strings.TrimSpace(password), PLAYERS) {
			// Handle the connection in a new goroutine
			_, err := conn.Write([]byte("successful"))
			checkError(err)
			CONNECTIONS[username] = conn
			fmt.Println(CONNECTIONS)

			initialCoord, err := infoReader.ReadString('\n')
			checkError(err)

			fmt.Println(string(initialCoord))

			go HandleConnection(conn)
		} else {
			conn.Write([]byte("failed"))
		}

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
