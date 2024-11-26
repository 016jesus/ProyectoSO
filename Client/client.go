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
	if len(os.Args) > 4 || len(os.Args) < 3 {
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


	conn, err := net.Dial("tcp4", ip + ":" + port)
	if err != nil {
		fmt.Println("Error conectando al servidor:", err)
		return
	}


	networkWriter := bufio.NewWriter(conn)
	networkReader := bufio.NewReader(conn)
	attemps, err := networkReader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo intentos de login:", err)
		return
	}
	attempts, err = strconv.Atoi(strings.Trim(attemps, "\n"))
	if err != nil {
		fmt.Println("Error convirtiendo intentos de login:", err)
		return
	}
	var username, password string
	var login bool = false
	buffer := bufio.NewReader(os.Stdin)
	for i := 0; i < attempts; i++ {
		if i != 0{
			fmt.Println("Verifique sus credenciales. Intentos restantes:", attempts - i)
		}
		fmt.Print("Login as: ")
		username, _ = buffer.ReadString('\n')
		fmt.Print("Password: ")
		password, _ = buffer.ReadString('\n')
		username = strings.TrimSpace(strings.Trim(username, "\n"))
		password = strings.TrimSpace(strings.Trim(password, "\n"))

		password = helpers.Encrypt(password)
		credentials := username + ":" + password
		
		networkWriter.WriteString(credentials + "\n")

		err = networkWriter.Flush()
		if err != nil {
			fmt.Println("Error enviando credenciales al servidor:", err)
			return
		}
		Reader := bufio.NewReader(conn)
		response, err := Reader.ReadString('\n')

		response = strings.Trim(response, "\n")
		if err != nil{
			fmt.Printf("Error leyendo respuesta del servidor: %v\n", err)
			continue
		}

		if response == "LOGIN_OK" {
			login = true
			break
		}
}

	if login {
		//enviar intervalo de tiempo
		networkWriter.WriteString(seconds + "\n")
		networkWriter.Flush()
		//ceder acceso a ejecucion de comandos
		helpers.ClientTCP(&conn)
	} else {
		fmt.Println("Demasiados intentos fallidos. Saliendo...")
		return
	}

}
