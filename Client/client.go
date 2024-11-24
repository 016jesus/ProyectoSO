package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"proyectoso/helpers"
	"strconv"
	"strings"
)

func main() {
	var seconds string
	if len(os.Args) > 4 && len(os.Args) < 3 {
		fmt.Println("Uso: ./client <server_ip> <port> <seconds>")
		return
	} else if len(os.Args) == 4{
		seconds = os.Args[3]
	} else {
		seconds = "5"
	}

	//parametros de ejecucion
	ip := os.Args[1]
	port := os.Args[2]


	//numero de intentos de login
	var attempts int
	output, err := helpers.ReadConfig("/etc/serverOper/serverOper.conf", "attempts")
	if err != nil{
		attempts = 3
	}else{
		attempts, _ = strconv.Atoi(output)
	}

	conn, err := net.Dial("tcp4", ip + ":" + port)
	if err != nil {
		fmt.Println("Error conectando al servidor:", err)
		return
	}
	defer conn.Close()

	localReader := bufio.NewReader(os.Stdin)
	networkWriter := bufio.NewWriter(conn)
	networkReader := bufio.NewReader(conn)
	for i := 0; i <= attempts; i++ {
		if i != 0{
			fmt.Println("Verifique sus credenciales. Intentos restantes:", attempts - i)
		}
		fmt.Print("Login as: ")
		username, _ := localReader.ReadString('\n')
		username = strings.Trim(username, "\n")
		fmt.Print("Password: ")
		password, _ := localReader.ReadString('\n')
		password = strings.Trim(password, "\n")
		password = helpers.Encrypt(password)
		credentials := username + ":" + password
		
		networkWriter.WriteString(credentials + "\n")

		if networkWriter.Flush() != nil {
			fmt.Println("Error sending credentials:", err)
			return
		}

		response, err := networkReader.ReadString('\n')
		response = strings.Trim(response, "\n")
		if err != nil{
			fmt.Printf("Error leyendo respuesta del servidor: %v\n", err)
			continue
		}
		fmt.Print(response)
		if response == "LOGIN_OK" {
			//enviar intervalo de tiempo
			networkWriter.WriteString(seconds + "\n")
			//ceder acceso a ejecucion de comandos
			helpers.ClientTCP(&conn)
		}
	}
}
