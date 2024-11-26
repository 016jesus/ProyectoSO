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
    var err error
    messenger := bufio.NewWriter(*socket)
    remoteReader := bufio.NewReader(*socket)

    scanner := bufio.NewScanner(os.Stdin)

    for {
        symbol, _ := remoteReader.ReadString('\n')
        fmt.Print(strings.TrimSpace(symbol) + " ")

        if !scanner.Scan() {
            fmt.Println("Error leyendo la entrada del usuario")
            break
        }
        comando := strings.TrimSpace(scanner.Text())

        if comando == "bye" {
            messenger.WriteString("bye\n")
            err := messenger.Flush()
            if err != nil {
                fmt.Printf("Error cerrando conexión %v", err)
            }
            fmt.Println("Cerrando conexión...")
            return
        } else {
            // Enviar comando al servidor
            messenger.WriteString(comando + "\n")
            messenger.Flush()

            // obtener la respuesta del servidor y mostrar en consola
            response, err = getOutput(remoteReader)
            if err != nil {
                fmt.Println(response.String())
                break;
            }
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

    go GetSystemReports(intervalo)

    for {
        symbol := getSystemSymbol()
        messenger.WriteString(symbol + "\n")
        messenger.Flush()
        command, err := remoteReader.ReadString('\n')
        if err != nil {
            fmt.Println("Error leyendo del socket:", err)
            break
        }
        command = strings.TrimSpace(command)
        if command == "" {
            continue
        }
      

        switch command {
        case "bye":
            fmt.Println("Shell cerrado por el cliente.")
            return
        case "report":
            control.RLock()
            reporte := reporteActual
            control.RUnlock()
            messenger.WriteString(reporte + "\n")
        case "report -r":
            messenger.WriteString(Ram + "\n")
        case "report -c":
            messenger.WriteString(Cpu + "\n")
        case "report -d":
            messenger.WriteString(Disk + "\n")
        default:
            output := exec.Command("/bin/sh", "-c", command)
            executedOutput, err := output.CombinedOutput()
            if err != nil {
                fmt.Print(err)
                messenger.WriteString("Nada apropiado\n")
            }else if len(executedOutput) == 0 {
                messenger.WriteString("Nada apropiado\n")
            } else {
                messenger.WriteString(string(executedOutput) + "\n")
            }
        }
        messenger.Flush()
    }
}




func getSystemSymbol() string {
	usr, _ := user.Current()
	host, _ := os.Hostname()
	return fmt.Sprintf("%s@%s:~$", usr.Username, host)
}