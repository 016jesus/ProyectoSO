package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
localReader := bufio.NewReader(os.Stdin)

comando, _ := localReader.ReadString('\n')
fmt.Print(comando)


}