package helpers

import (
	"fmt"
	"crypto/sha256"
)

//Comprueba las credenciales de usuario
func ValidarLogin(array []string, credb map[string]string) bool {
	user := array[0]
	passwd, ok := credb[user]
	if ok && passwd == array[1] {
			return true
	}
	return false
}

//retorna la passwd hasheada y en hexadecimal
func Encrypt(passwd string)string{
	return fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))
}