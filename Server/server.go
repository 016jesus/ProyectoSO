package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"proyectoso/helpers"
	"strconv"
	"strings"
	"time"
)

func main() {

	mapasswd := helpers.ReadCredentials("/etc/serverOper/users.db")

	fmt.Println("**************************************")
	fmt.Println("**** ServerOper - JGD, DAMT, MJMZ ****")
	fmt.Println("**************************************")
	attempts, _:= helpers.ReadConfig("/etc/serverOper/serverOper.conf", "attempts")
	attempts = strings.Trim(attempts, "\n")
	//preparar conexion
	direccionTCP, err := net.ResolveTCPAddr("tcp4", ":2024")
	if err != nil {
		helpers.WriteLog("Error resolviendo dirección:" + err.Error())
		fmt.Println("Error resolviendo dirección:", err)
		return
	}
	socket, err := net.ListenTCP("tcp", direccionTCP)
	if err != nil {
		helpers.WriteLog("Error iniciando servidor:" + err.Error())
		fmt.Println("Error iniciando servidor:", err)
		return
	}
	defer socket.Close()
	fmt.Println("Servidor escuchando en ", socket.Addr().String())
	helpers.WriteLog("Servidor abierto en " + socket.Addr().String())
	seconds:= 5
	//esperar conexiones
	for{
		fmt.Println("Esperando conexion...")
		conn , err := socket.Accept()
		if err != nil {
			helpers.WriteLog(err.Error())
			log.Fatal(err)
		}
			fmt.Println("Conexion establecida desde: ", conn.RemoteAddr().String())
			
		
			
			messenger := bufio.NewWriter(conn)
			messenger.WriteString(attempts + "\n")
			if messenger.Flush() != nil {
				log.Fatal("Error enviando intentos")
			}
			
			limit, _ := strconv.Atoi(attempts)
			fmt.Print("Intentos: ", limit)
			for i := 0; i <= limit; i++ {
			buffer := bufio.NewReader(conn)
			credentials:= helpers.ReceiveCredentials(buffer)
			
			if helpers.ValidarLogin(credentials, mapasswd) {
				
				// ceder el control a la funcion de conexion
				_, err := messenger.WriteString("LOGIN_OK\n")
				if err != nil {
					log.Fatal("Error enviando mensaje de login:", err)
				}
				err = messenger.Flush()
				if err != nil {
					log.Fatal("Error enviando mensaje de login:", err)
				}
				mensaje := "Cliente " + credentials[0] + " autenticado en: " + conn.RemoteAddr().String()
				helpers.WriteLog(mensaje)
				// recibir intervalo de tiempo
				fmt.Println(mensaje)
				intervalo, err := buffer.ReadString('\n')
				if err != nil {
					log.Fatal("Error leyendo intervalo de tiempo:", err)
				}
				seconds, err = strconv.Atoi(strings.Trim(intervalo, "\n"))
				if err != nil {
					log.Fatal("Error convirtiendo intervalo de tiempo:", err)
				}
				// ceder control a la funcion de ejecucion de comandos
				break
				
			}else{
				messenger.WriteString("LOGIN_FAIL\n")
				mensaje := "Error autenticando desde: " + conn.RemoteAddr().String()
				fmt.Println(mensaje)
				helpers.WriteLog(mensaje)
				messenger.Flush()
				
			}
		}
		go helpers.ServerTCP(&conn, time.Duration(seconds))
	}
}
