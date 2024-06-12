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
	"time"

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
var MAP = make(map[string]string)

var PLAYERS []Player //for authentication
var ROWS, COLS = 10, 18
var BOARD = make([][]string, ROWS)
var CONNECTIONS = make(map[string]net.Conn)
var despawnQueues []string
var POKEMON_LOCATIONS = make(map[string]string)
var PLAYER_LOCATIONS = make(map[string]string)

// ////////////////////////////////////////////////////////////////////////////////////
// Load players.json to PLAYERS
func loadPlayers(filename string) []Player {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	var players []Player
	err = json.Unmarshal(bytes, &players)
	if err != nil {
		return nil
	}

	return players
}

func loadPokemons(filename string) []Pokemon {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	var pokemons []Pokemon
	err = json.Unmarshal(bytes, &pokemons)
	if err != nil {
		return nil
	}

	return pokemons
}
func verifyPlayer(username, password string, players []Player) bool {
	for _, user := range players {
		if user.Username == username {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			return err == nil
		}
	}
	return false
}
func generateRandomPokemons(num int) map[string]string {
	pokemonLocations := make(map[string]string)
	for range num {
		for {
			spawnX := rand.Intn(ROWS)
			spawnY := rand.Intn(COLS)
			if BOARD[spawnX][spawnY] == "" {
				pokemonID := POKEMONS[rand.Intn(len(POKEMONS))].ID

				BOARD[spawnX][spawnY] = pokemonID

				despawnQueues = append(despawnQueues, strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY))
				pokemonLocations[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = pokemonID
				POKEMON_LOCATIONS[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = pokemonID
				break
			}

		}

	}
	return pokemonLocations
}

// ////////////////////////////////////////////////////////////////////////////////////
func handlePokemons() {
	spawnTicker1min := time.NewTicker(1 * time.Minute)
	despawnTicker5min := time.NewTicker(5 * time.Minute)
	NUMBERTOPROCESS := 3
	for {
		select {
		case <-spawnTicker1min.C:

			newPokemonLocations, err := json.Marshal(generateRandomPokemons(NUMBERTOPROCESS))
			checkError(err)

			for _, tcpConn := range CONNECTIONS {

				tcpConn.Write([]byte(newPokemonLocations))

			}
		case <-despawnTicker5min.C:

		}
	}

}
func battle(player1 string, player2 string) {

}

// HandleConnection handles incoming client connections
func HandleConnection(conn net.Conn) {

	defer conn.Close()

	// Create a reader to read data from the connection
	reader := bufio.NewReader(conn)
	for {
		battleStatus := false
		// Read data from the connection
		player_coord, err := reader.ReadString('\n')

		//DISCONNECTED
		if err != nil {
			for name, connection := range CONNECTIONS {
				//tim connection hien tai
				if connection == conn {
					// tim ten
					for location, player := range PLAYER_LOCATIONS {
						if player == name {
							delete(PLAYER_LOCATIONS, location)
						}
					}
					delete(CONNECTIONS, name)
					quit := make(map[string]string)
					quit[strings.TrimSpace(name)] = "quit"
					sentQuit, _ := json.Marshal(quit)
					for _, connection := range CONNECTIONS {

						connection.Write([]byte(sentQuit))
					}
					fmt.Println(strings.TrimSpace(name) + " is disconnected")
					break

				}

			}
			return

		}

		if strings.Split(player_coord, "-")[0] == "battle" {

		}

		// iterate through connections
		for name, connection := range CONNECTIONS {
			// find matching name
			if connection == conn {
				//find the name corresponding to connection address
				for playerLocation, player := range PLAYER_LOCATIONS {

					if player == name {

						//if player hits pokemon location
						if pokemonIndex, exists := POKEMON_LOCATIONS[strings.TrimSpace(player_coord)]; exists {
							fmt.Println("CATCHING...")
							//send pokemon id to the player through "conn"
							// catched: "index"
							catched := make(map[string]string)
							catched[strings.TrimSpace(name)] = pokemonIndex

							sentCatched, _ := json.Marshal(catched)
							conn.Write(sentCatched)
							//remove location of the pokemon from BOARD and POKEMON_LOCATIONS
							pokemonX, _ := strconv.Atoi(strings.Split(playerLocation, "-")[0])
							pokemonY, _ := strconv.Atoi(strings.Split(playerLocation, "-")[1])
							BOARD[pokemonX][pokemonY] = ""
							delete(POKEMON_LOCATIONS, strings.TrimSpace(player_coord))
							fmt.Println(POKEMON_LOCATIONS)
							for _, tcpConn := range CONNECTIONS {
								if tcpConn != conn {
									pokemonGone := make(map[string]string)
									pokemonGone[strings.TrimSpace(player_coord)] = ""
									sentPokemonGone, _ := json.Marshal(pokemonGone)
									tcpConn.Write([]byte(sentPokemonGone))
								}
							}
						} else {
							for enemy_location, enemy_name := range PLAYER_LOCATIONS {
								if enemy_location == strings.TrimSpace(player_coord) && enemy_name != name {
									//BATTLE
									fmt.Println("battle")
									//send battle noti to player
									battleInfo := make(map[string]string)
									battleInfo["battle"] = enemy_name
									sentBattleInfo, _ := json.Marshal(battleInfo)
									conn.Write([]byte(sentBattleInfo))

									//send to the enemy
									battledInfo := make(map[string]string)
									battledInfo["battle"] = name
									sentBattledInfo, _ := json.Marshal(battledInfo)
									CONNECTIONS[enemy_name].Write([]byte(sentBattledInfo))

									battleStatus = true
									// start new routine for
									// go battle(name, enemy_name)
								}
							}
						}

						//remove previous location to make animation
						delete(PLAYER_LOCATIONS, playerLocation)
					}
				}
				if !battleStatus {
					PLAYER_LOCATIONS[player_coord] = name
				}
			}
		}

		// Echo the message back to the client

		for _, tcpConn := range CONNECTIONS {
			// if conn != tcpConn {
			sentPLAYER_LOCATIONS, _ := json.Marshal(PLAYER_LOCATIONS)
			tcpConn.Write([]byte(sentPLAYER_LOCATIONS))

			// }
		}
	}
}
func handle(conn net.Conn) {
	// Verfify username and password
	infoReader := bufio.NewReader(conn)

	// Get username
	username, err := infoReader.ReadString('\n')
	checkError(err)

	// Get password
	password, err := infoReader.ReadString('\n')
	checkError(err)

	if verifyPlayer(strings.TrimSpace(username), strings.TrimSpace(password), PLAYERS) {
		// Handle the connection in a new goroutine
		_, err := conn.Write([]byte("successful"))
		initialPokemons := ""
		for i := range 3 {
			initialPokemons += strconv.Itoa(rand.Intn(len(POKEMONS)))
			if i < 2 {
				initialPokemons += "-"
			}
		}
		conn.Write([]byte(initialPokemons))
		checkError(err)
		time.Sleep(22 * time.Second)
		CONNECTIONS[username] = conn
		fmt.Println(CONNECTIONS)

		//send pokemon locations
		sentPOKEMON_LOCATIONS, _ := json.Marshal(POKEMON_LOCATIONS)
		conn.Write([]byte(sentPOKEMON_LOCATIONS))
		for {
			playerX := rand.Intn(ROWS)
			playerY := rand.Intn(COLS)
			if BOARD[playerX][playerY] == "" {
				BOARD[playerX][playerY] = username
				PLAYER_LOCATIONS[strconv.Itoa(playerX)+"-"+strconv.Itoa(playerY)] = username

				// ini_coord := make(map[string]string)
				// ini_coord[strconv.Itoa(playerX)+"-"+strconv.Itoa(playerY)] = username

				// //send initial player's location
				// sentInital_Coord, _ := json.Marshal(ini_coord)
				// conn.Write([]byte(sentInital_Coord))
				break
			}
		}

		//send player locations
		sentPLAYER_LOCATIONS, _ := json.Marshal(PLAYER_LOCATIONS)
		for _, connection := range CONNECTIONS {
			connection.Write([]byte(sentPLAYER_LOCATIONS))
		}

		go HandleConnection(conn)
	} else {
		conn.Write([]byte("failed"))
	}
}
func main() {
	// Set up BOARD
	for i := range BOARD {
		BOARD[i] = make([]string, COLS)
	}
	// Load pokemon.json to array POKEMONS
	POKEMONS = loadPokemons("pokedex.json")

	//Load players.json to array PLAYERS
	PLAYERS = loadPlayers("players.json")
	generateRandomPokemons(3)
	fmt.Println(POKEMON_LOCATIONS)
	go handlePokemons()

	// Start listening for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	/// spawn pokem
	for {
		// Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handle(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
