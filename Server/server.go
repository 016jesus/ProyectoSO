package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"proyectoso/helpers"
)


func main() {
	mapasswd := helpers.ReadCredentials("/etc/serverOper/users.db")



	fmt.Println("**************************************")
	fmt.Println("**** ServerOper - JGD, DAMT, MJMZ ****")
	fmt.Println("**************************************")



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
	for{
		conn , err := socket.Accept()
		if err != nil {
			helpers.WriteLog(err.Error())
			log.Fatal(err)
		}
		buffer := bufio.NewReader(conn)
		credentials:= helpers.ReceiveCredentials(buffer)
		messenger := bufio.NewWriter(conn)
		if(helpers.ValidarLogin(credentials, mapasswd)){
			//ceder el control a la funcion de conexion
			messenger.WriteString("SUCCESSFUL_LOGIN\n")
			messenger.Flush()
			helpers.WriteLog("Cliente "+ credentials[0] + " autenticado en: " + conn.RemoteAddr().Network())
			go helpers.ServerTCP(&conn)

		}
	}
}