package helpers

import (
	"fmt"
	"crypto/sha256"
)

//Comprueba las credenciales de usuario
func ValidarLogin(array []string, credb map[string]string) bool {
	user := array[0]
	if passwd, ok := credb[user]; ok {
		if passwd == array[1] {
			return true
		}
	}
	return true
}

//retorna la passwd hasheada y en hexadecimal
func Encrypt(passwd string)string{
	return fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))
}