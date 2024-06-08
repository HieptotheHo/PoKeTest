package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

var ROWS, COLS = 10, 18
var BOARD = make([][]string, ROWS)
var USERNAME = ""

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
				if cell == USERNAME {
					fmt.Printf("| %s ", "☠")
				} else {
					fmt.Printf("| %s ", "?")
				}
			}
		}
		fmt.Println("|")
	}
	// Print the final horizontal line after all rows
	fmt.Println(horizontalLine(len(board[0])))
}

// read from the server and print to console
func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// Read server's response
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}

		// Print the message from the server
		// fmt.Print("Server response: " + message)
		drawBoard(BOARD)
	}
}

func main() {
	//INITIALIZING MAP
	for i := range BOARD {
		BOARD[i] = make([]string, COLS)
	}

	X := rand.Intn(ROWS)
	Y := rand.Intn(COLS)

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
		USERNAME = username
		BOARD[X][Y] = username
		_, err := conn.Write([]byte(strconv.Itoa(X) + "-" + strconv.Itoa(Y) + "\n"))
		checkError(err)
		go readFromServer(conn)
		fmt.Println("MAIN GAME:")
		// for scanner.Scan() {
		// 	message := scanner.Text()
		// 	if message == "exit" {
		// 		fmt.Println("Exiting...")
		// 		break
		// 	}

		// 	// Send the message to the server
		// 	_, err := conn.Write([]byte(message + "\n"))
		// 	if err != nil {
		// 		fmt.Println("Error writing to server:", err)
		// 		break
		// 	}
		// }
		// if err := scanner.Err(); err != nil {
		// 	fmt.Println("Error reading from stdin:", err)
		// }
		// Initialize the keyboard
		if err := keyboard.Open(); err != nil {
			fmt.Println("Failed to open keyboard:", err)
		}
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println("Error reading key:", err)
				continue
			}

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

			drawBoard(BOARD)
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}