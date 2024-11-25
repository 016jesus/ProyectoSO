package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"log"
)

func ReadCredentials(ruta string) map[string]string {
	// Leer el archivo
	file, err := os.Open(ruta)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return nil

	}
	defer file.Close()
	buffer := bufio.NewReader(file)


	credentialsMap := make(map[string]string)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }
		credentials := strings.Split(line, ":")
		if len(credentials) == 2 {
			credentialsMap[credentials[0]] = credentials[1]
		} else {
			log.Fatal("Error en el formato del archivo de credenciales")
		}
	}
	return credentialsMap
}


func ReadConfig(ruta string, config string) (string, error){
	file, err := os.Open(ruta)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buffer := bufio.NewReader(file)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
            if err == io.EOF {
                break
            }
			return "", err
        }
		credentials := strings.Split(line, "=")
		if credentials[0] == config {
			return credentials[0], nil
		} else {
			continue	
		}
	}
	return "", fmt.Errorf("%s","No se encontro la configuracion suministrada")
}