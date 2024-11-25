package helpers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
	"time"
)

/*
*
* Variables globales
*
*/
var reporteActual string
var control sync.RWMutex



/*
*
* Funcion de conexion del lado del cliente
*
*/
func ClientTCP(socket *net.Conn) {
	defer (*socket).Close()
	var response strings.Builder
	messenger := bufio.NewWriter(*socket)
	remoteReader := bufio.NewReader(*socket)

	for {
		symbol, _ := remoteReader.ReadString('\n')
		fmt.Print(strings.TrimSpace(symbol) + " ")
		
		localReader := bufio.NewReader(os.Stdin)
		comando, _ := localReader.ReadString('\n')
		comando = strings.TrimSpace(comando)
		if comando == "bye" {
			messenger.WriteString("bye\n")
			err := messenger.Flush()
			if  err != nil {
				fmt.Printf("Error cerrando conexión %v", err)
			}
			fmt.Println("Cerrando conexión...")
			return
		} else{
		// Enviar comando al servidor
		messenger.WriteString(comando + "\n")
		messenger.Flush()

		// obtener la respuesta del servidor y mostrar en consola
		response = getOutput(remoteReader)
		fmt.Println(response.String())
		}
}

}




/*
*
* Funcion de conexion del lado del servidor
*
*/


func ServerTCP(socket *net.Conn, intervalo time.Duration) {

	messenger := bufio.NewWriter(*socket)
	remoteReader := bufio.NewReader(*socket)

	// Iniciar generación de reportes periódicos

	go GetSystemReports(intervalo)

	for {
		// Leer el símbolo del sistema
		symbol := getSystemSymbol()
		messenger.WriteString(symbol + "\n")
		//enviar simbolo del sistema al cliente y leer comando
		messenger.Flush()
		command, err := remoteReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error leyendo del socket:", err)
			break
		}
		// limpiar el comando
		command = strings.TrimRight(command, "\r\n")
		// Verificar si hay comandos
		if command == "" {
			continue
		}
		if command == "bye" {
			fmt.Println("Shell cerrado por el cliente.")
			break
		}

		if command == "report" {
			// Leer el último reporte generado
			control.RLock()
			reporte := reporteActual
			control.RUnlock()
			fmt.Print("Enviando reporte al cliente...", reporte)
			_, _ = messenger.WriteString(reporte + "\n")
			messenger.Flush()
			continue
		}

		output := exec.Command("/bin/sh", "-c", command)
		executedOutput, err := output.CombinedOutput()
		if err != nil {
			salidaError := fmt.Sprintf("Error ejecutando el comando: %s\n", err.Error())
			_, _ = messenger.WriteString(salidaError)
			messenger.Flush()
			continue
		}

		// Enviar salida del comando al cliente
		commandOutput := string(executedOutput)
		messenger.WriteString(commandOutput + "\n")
		messenger.Flush()

		fmt.Printf("\nEnviado al cliente %s: %s\n", (*socket).RemoteAddr().String(), commandOutput)
	}
}




func getSystemSymbol() string {
	usr, _ := user.Current()
	host, _ := os.Hostname()
	return fmt.Sprintf("%s@%s:~$ ", usr.Username, host)
}