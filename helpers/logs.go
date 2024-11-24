package helpers

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)


func ReceiveCredentials(buffer *bufio.Reader) []string {
	line, err := buffer.ReadString('\n')
	if err != nil {
		WriteLog("Error leyendo la línea: " + err.Error())
		log.Fatalf("Error leyendo la línea: %v", err)
	}

	credentials := strings.SplitN(line, ":", 2)
	if len(credentials) != 2 {
		WriteLog("Formato de credenciales inválido")
		log.Fatalf("Formato de credenciales inválido")
	}

	return credentials
}

func WriteLog(mensaje string) {
	logFile, err := os.OpenFile("/var/log/serverOper.log", os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Fatalf("Error abriendo el archivo de logs: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	timestamp := time.Now().Format("2024-01-21 15:04:05")
	logger.Printf("%s %s", timestamp, mensaje)
}
