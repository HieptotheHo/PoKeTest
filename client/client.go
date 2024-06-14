package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"

	gim "github.com/ozankasikci/go-image-merge"
)

var ROWS, COLS = 10, 18
var BOARD = make([][]string, ROWS)
var USERNAME = ""
var X int
var Y int

var ENEMIES = make(map[string]string)

type Pokemon struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Types []string          `json:"types"`
	Stats map[string]string `json:"stats"`
	Exp   string            `json:"exp"`
}

var POKEMONS []Pokemon

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
func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
func drawTitle() {
	fmt.Println("                                  ,'\\")
	fmt.Println("    _.----.        ____         ,'  _\\   ___    ___     ____")
	fmt.Println("_,-'       `.     |    |  /`.   \\,-'    |   \\  /   |   |    \\  |`.")
	fmt.Println("\\      __    \\    '-.  | /   `.  ___    |    \\/    |   '-.   \\ |  |")
	fmt.Println(" \\.    \\ \\   |  __  |  |/    ,','_  `.  |          | __  |    \\|  |")
	fmt.Println("   \\    \\/   /,' _`.|      ,' / / / /   |          ,' _`.|     |  |")
	fmt.Println("    \\     ,-'/  / \\ \\    ,'   | \\/ / ,`.|         /  / \\ \\  |     |")
	fmt.Println("     \\    \\ |   \\_/  |   `-.  \\    `'  /|  |    ||   \\_/  | |\\    |")
	fmt.Println("      \\    \\ \\      /       `-.`.___,-' |  |\\  /| \\      /  | |   |")
	fmt.Println("       \\    \\ `.__,'|  |`-._    `|      |__| \\/ |  `.__,'|  | |   |")
	fmt.Println("        \\_.-'       |__|    `-._ |              '-.|     '-.| |   |")
	fmt.Println("                                `'                            '-._|")
}
func drawBoard(board [][]string) {

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	drawTitle()
	// Function to generate a horizontal line
	horizontalLine := func(length int) string {
		return "+" + strings.Repeat("---+", length)
	}

	for _, row := range board {
		// Print horizontal line before each row
		fmt.Println(horizontalLine(len(row)))

		// Print cell values or empty spaces ☺ ☻ ☠
		for _, cell := range row {
			if cell == "" {
				fmt.Print("|   ")
			} else {
				if isNumber(cell) {
					fmt.Printf("| %s ", "?")
				} else {

					if cell == USERNAME {
						fmt.Printf("| %s ", "☻")
					} else {
						fmt.Printf("| %s ", "☠")
					}
				}

			}
		}
		fmt.Println("|")
	}
	// Print the final horizontal line after all rows
	fmt.Println(horizontalLine(len(board[0])))
}
func drawCongrats() {
	fmt.Println("░█████╗░░█████╗░███╗░░██╗░██████╗░██████╗░░█████╗░████████╗░██████╗")
	fmt.Println("██╔══██╗██╔══██╗████╗░██║██╔════╝░██╔══██╗██╔══██╗╚══██╔══╝██╔════╝")
	fmt.Println("██║░░╚═╝██║░░██║██╔██╗██║██║░░██╗░██████╔╝███████║░░░██║░░░╚█████╗░")
	fmt.Println("██║░░██╗██║░░██║██║╚████║██║░░╚██╗██╔══██╗██╔══██║░░░██║░░░░╚═══██╗")
	fmt.Println("╚█████╔╝╚█████╔╝██║░╚███║╚██████╔╝██║░░██║██║░░██║░░░██║░░░██████╔╝")
	fmt.Println("░╚════╝░░╚════╝░╚═╝░░╚══╝░╚═════╝░╚═╝░░╚═╝╚═╝░░╚═╝░░░╚═╝░░░╚═════╝░")
}
func drawStats(pokemon Pokemon) {
	fmt.Println(pokemon.Name)
	fmt.Print("Types: ")
	for i := 0; i < len(pokemon.Types); i++ {
		fmt.Print(pokemon.Types[i] + " ")
	}
	fmt.Println()

	fmt.Print("HP:              ")
	hp, _ := strconv.Atoi(pokemon.Stats["HP"])
	for i := 0; i < hp; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("ATTACK:          ")
	attk, _ := strconv.Atoi(pokemon.Stats["Attack"])
	for i := 0; i < attk; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("SPECIAL ATTACK:  ")
	sp_atk, _ := strconv.Atoi(pokemon.Stats["Sp Atk"])
	for i := 0; i < sp_atk; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("DEFENSE:         ")
	df, _ := strconv.Atoi(pokemon.Stats["Defense"])
	for i := 0; i < df; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("SPECIAL DEFENSE: ")
	sp_def, _ := strconv.Atoi(pokemon.Stats["Sp Def"])
	for i := 0; i < sp_def; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("SPEED:           ")
	speed, _ := strconv.Atoi(pokemon.Stats["Speed"])
	for i := 0; i < speed; i++ {
		fmt.Print("█")
	}
	fmt.Println()
	fmt.Println()
}

var pokeBalls []Pokemon
var DRAWBOARD_SIGNAL = true

func showNewPoKemon(pokemon Pokemon) {
	for {

		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		drawCongrats()
		drawStats(pokemon)
		cmd = exec.Command("cmd", "/c", "image2ascii.exe -f .\\pokemon_images\\pokemon_"+pokemon.ID+".png -r 0.7")
		cmd.Stdout = os.Stdout
		cmd.Run()
		pokeBalls = append(pokeBalls, pokemon)
		fmt.Println(pokeBalls)
		time.Sleep(5 * time.Second)
		DRAWBOARD_SIGNAL = true
		drawBoard(BOARD)
		break

	}
}

func displayDeck() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("LET'S BATTLE!")
	for i := 0; i < len(pokeBalls); i++ {
		str := "dh12788$!@1ufiu1!@%!F"

		shuff := []rune(str)

		// Shuffling the string
		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})

		if i%3 == 0 {
			names := pokeBalls[i].Name
			if i+1 < len(pokeBalls) {
				names += "\t\t" + pokeBalls[i+1].Name
			}
			if i+2 < len(pokeBalls) {
				names += "\t\t" + pokeBalls[i+2].Name
			}
			fmt.Println(names)
			if i+2 < len(pokeBalls) {
				// Shuffling the string
				rand.Shuffle(len(shuff), func(i, j int) {
					shuff[i], shuff[j] = shuff[j], shuff[i]
				})
				grids := []*gim.Grid{
					{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i].ID + ".png"},
					{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i+1].ID + ".png"},
					{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i+2].ID + ".png"},
				}
				rgba, _ := gim.New(grids, 3, 1).Merge()
				file, _ := os.Create("./deck/" + string(shuff) + ".png")
				jpeg.Encode(file, rgba, &jpeg.Options{Quality: 80})
				png.Encode(file, rgba)

				cmd = exec.Command("cmd", "/c", "image2ascii.exe -f .\\deck\\"+string(shuff)+".png -r 0.3")
				cmd.Stdout = os.Stdout
				cmd.Run()
			} else {
				if i+1 < len(pokeBalls) {
					// Shuffling the string
					rand.Shuffle(len(shuff), func(i, j int) {
						shuff[i], shuff[j] = shuff[j], shuff[i]
					})
					grids := []*gim.Grid{
						{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i].ID + ".png"},
						{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i+1].ID + ".png"},
					}
					rgba, _ := gim.New(grids, 3, 1).Merge()
					file, _ := os.Create("./deck/" + string(shuff) + ".png")
					jpeg.Encode(file, rgba, &jpeg.Options{Quality: 80})
					png.Encode(file, rgba)

					cmd = exec.Command("cmd", "/c", "image2ascii.exe -f .\\deck\\"+string(shuff)+".png -r 0.3")
					cmd.Stdout = os.Stdout
					cmd.Run()
				} else {
					// Shuffling the string
					rand.Shuffle(len(shuff), func(i, j int) {
						shuff[i], shuff[j] = shuff[j], shuff[i]
					})
					grids := []*gim.Grid{
						{ImageFilePath: "./pokemon_images/pokemon_" + pokeBalls[i].ID + ".png"},
					}
					rgba, _ := gim.New(grids, 3, 1).Merge()
					file, _ := os.Create("./deck/" + string(shuff) + ".png")
					jpeg.Encode(file, rgba, &jpeg.Options{Quality: 80})
					png.Encode(file, rgba)

					cmd = exec.Command("cmd", "/c", "image2ascii.exe -f .\\deck\\"+string(shuff)+".png -r 0.3")
					cmd.Stdout = os.Stdout
					cmd.Run()
				}
			}
		}
	}
}
func battle() {
	displayDeck()
	for {

		// cmd = exec.Command("cmd", "/c", "image2ascii.exe -f ./path.png -r 0.5")
		// cmd.Stdout = os.Stdout
		// cmd.Run()
	}
}

// read from the server and print to console
func readFromServer(conn net.Conn) {
	// reader := bufio.NewReader(conn)

	for {
		buf := make([]byte, 1024)
		// Read server's response
		// message, err := reader.ReadString('\n')

		n, err := conn.Read(buf)
		checkError(err)

		// Print the message from the server
		var locations map[string]string

		err = json.Unmarshal(buf[:n], &locations)
		checkError(err)
		// fmt.Println("Server response: ", locations)

		for location, id := range locations {
			location = strings.TrimSpace(location)
			id = strings.TrimSpace(id)
			//battle
			if location == "battle" {

				if id == USERNAME {
					//MY TURN
					fmt.Println("hello")
				} else {
					displayDeck()
					//enter pokemons for battling
					var chosenPokemons []Pokemon
					var pokemonName string

					for {
						fmt.Println("Choose your Pokemons")
						fmt.Print("Name: ")
						fmt.Scanln(&pokemonName)

						for pokeIndex := range pokeBalls {
							if pokeBalls[pokeIndex].Name == pokemonName {
								fmt.Println("exists")
								chosenPokemons = append(chosenPokemons, pokeBalls[pokeIndex])
								conn.Write([]byte("battle-" + USERNAME + "-" + (pokeBalls[pokeIndex].ID + "\n")))
							}
						}
						if len(chosenPokemons) == 3 {
							cmd := exec.Command("cmd", "/c", "cls")
							cmd.Stdout = os.Stdout
							cmd.Run()
							break
						}
					}

					fmt.Println("Wait for your opponent to choose pokemons...")

					// Read server's response
					// message, err := reader.ReadString('\n')

					DRAWBOARD_SIGNAL = true
					drawBoard(BOARD)
					/////
				}

				//send to server your move in battling
				conn.Write([]byte("battle-" + strings.TrimSpace(id)))

				DRAWBOARD_SIGNAL = false
			} else {
				if location == USERNAME {
					//get pokemon
					if isNumber(id) {
						BOARD[X][Y] = ""
						index, _ := strconv.Atoi(id)
						newPokemon := POKEMONS[index]

						go showNewPoKemon(newPokemon)
						DRAWBOARD_SIGNAL = false
					}

				} else {
					//HANDLE DISCONNECTION
					if id == "quit" {
						fmt.Println(location + " is disconnected")
						for ene_location, enemy := range ENEMIES {
							if enemy == location {
								enemyX, _ := strconv.Atoi(strings.Split(ene_location, "-")[0])
								enemyY, _ := strconv.Atoi(strings.Split(ene_location, "-")[1])

								BOARD[enemyX][enemyY] = ""
								fmt.Println(ENEMIES)
								delete(ENEMIES, ene_location)

								break
							}
						}

					} else {
						//RECEIVE MOVEMENT FROM PLAYERS AND SPAWN COORD OF POKEMONS AS A "MAP"
						spawnX, _ := strconv.Atoi(strings.Split(location, "-")[0])
						spawnY, _ := strconv.Atoi(strings.Split(location, "-")[1])

						if isNumber(id) || id == "" {
							//pokemon
							BOARD[spawnX][spawnY] = id
						} else {

							//player: you and enemies
							if id == USERNAME {
								//you
								X = spawnX
								Y = spawnY
								BOARD[spawnX][spawnY] = USERNAME
							} else {
								//enemies
								for ene_location, enemy := range ENEMIES {
									if enemy == id {
										enemyX, _ := strconv.Atoi(strings.Split(ene_location, "-")[0])
										enemyY, _ := strconv.Atoi(strings.Split(ene_location, "-")[1])

										BOARD[enemyX][enemyY] = ""
										delete(ENEMIES, ene_location)

									}

								}
								ENEMIES[strconv.Itoa(spawnX)+"-"+strconv.Itoa(spawnY)] = id
								BOARD[spawnX][spawnY] = "enemy"

							}

						}
					}

				}
			}

		}
		if DRAWBOARD_SIGNAL {
			drawBoard(BOARD)
		}

	}
}

func main() {
	//INITIALIZING MAP
	for i := range BOARD {
		BOARD[i] = make([]string, COLS)
	}
	POKEMONS = loadPokemons("pokedex.json")
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start a goroutine to handle server responses

	// Read from stdin and send to server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Username: ")
	scanner.Scan()

	username := scanner.Text()
	_, err = conn.Write([]byte(username + "\n"))
	checkError(err)

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()
	_, err = conn.Write([]byte(password + "\n"))
	checkError(err)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	checkError(err)

	if strings.TrimSpace(string(buffer[:n])) == "successful" {
		n, err = conn.Read(buffer)
		checkError(err)
		pokemonIndexes := strings.Split(strings.TrimSpace(string(buffer[:n])), "-")
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		cmd = exec.Command("cmd", "/c", "image2ascii.exe -f .\\title.png -r 0.3")
		cmd.Stdout = os.Stdout
		cmd.Run()
		time.Sleep(5 * time.Second)
		for _, pokemonIndex := range pokemonIndexes {
			index, _ := strconv.Atoi(pokemonIndex)
			showNewPoKemon(POKEMONS[index])
		}
		USERNAME = username
		go readFromServer(conn)
		fmt.Println("MAIN GAME:")
		// Initialize the keyboard
		if err := keyboard.Open(); err != nil {
			fmt.Println("Failed to open keyboard:", err)
		}
		for {
			_, key, err := keyboard.GetKey()
			checkError(err)

			switch key {

			case keyboard.KeyArrowUp:
				if X > 0 {
					BOARD[X][Y] = ""
					X--
					BOARD[X][Y] = USERNAME
					_, err := conn.Write([]byte(strconv.Itoa(X) + "-" + strconv.Itoa(Y) + "\n"))
					checkError(err)
				}

			case keyboard.KeyArrowDown:
				if X < ROWS-1 {
					BOARD[X][Y] = ""
					X++
					BOARD[X][Y] = USERNAME
					_, err := conn.Write([]byte(strconv.Itoa(X) + "-" + strconv.Itoa(Y) + "\n"))
					checkError(err)
				}

			case keyboard.KeyArrowLeft:
				if Y > 0 {
					BOARD[X][Y] = ""
					Y--
					BOARD[X][Y] = USERNAME
					_, err := conn.Write([]byte(strconv.Itoa(X) + "-" + strconv.Itoa(Y) + "\n"))
					checkError(err)
				}

			case keyboard.KeyArrowRight:

				if Y < COLS-1 {
					BOARD[X][Y] = ""
					Y++
					BOARD[X][Y] = USERNAME
					_, err := conn.Write([]byte(strconv.Itoa(X) + "-" + strconv.Itoa(Y) + "\n"))
					checkError(err)
				}

			case keyboard.KeyEsc:
				fmt.Println("Exiting...")
				return
			}
			// drawBoard(BOARD)
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err)
		os.Exit(1)
	}
}
