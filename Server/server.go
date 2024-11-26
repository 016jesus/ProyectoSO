package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "proyectoso/helpers"
    "strconv"
    "strings"
    "syscall"
    "time"
)

func main() {
    mapasswd := helpers.ReadCredentials("/etc/serverOper/users.db")

    fmt.Println("**************************************")
    fmt.Println("**** ServerOper - JGD, DAMT, MJMZ ****")
    fmt.Println("**************************************")
    attempts, _ := helpers.ReadConfig("/etc/serverOper/serverOper.conf", "attempts")
    attempts = strings.Trim(attempts, "\n")

    // preparar conexion
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

    // Manejar señales de interrupción para cerrar el servidor correctamente
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-signalChan
        fmt.Println("\nCerrando servidor...")
        socket.Close()
        os.Exit(0)
    }()

    // esperar conexiones
    for {
        fmt.Println("Esperando conexion...")
        conn, err := socket.Accept()
        if err != nil {
            helpers.WriteLog(err.Error())
            log.Fatal(err)
        }
        fmt.Println("Conexion establecida desde: ", conn.RemoteAddr().String())

        go handleConnection(conn, mapasswd, attempts)
    }
}

func handleConnection(conn net.Conn, mapasswd map[string]string, attempts string) {
    defer conn.Close()
    buffer := bufio.NewReader(conn)
    messenger := bufio.NewWriter(conn)
    messenger.WriteString(attempts + "\n")
    if messenger.Flush() != nil {
        log.Fatal("Error enviando intentos")
    }
    limit, _ := strconv.Atoi(attempts)
    for i := 0; i < limit; i++ {
        credentials := helpers.ReceiveCredentials(buffer)
        if helpers.ValidarLogin(credentials, mapasswd) {
            messenger.WriteString("LOGIN_OK\n")
            messenger.Flush()
            mensaje := "Cliente " + credentials[0] + " autenticado en: " + conn.RemoteAddr().String()
            helpers.WriteLog(mensaje)
            intervalo, _ := buffer.ReadString('\n')
            seconds, _ := strconv.Atoi(strings.Trim(intervalo, "\n"))
            go helpers.ServerTCP(&conn, time.Duration(seconds))
            return
        } else {
            messenger.WriteString("LOGIN_FAIL\n")
            messenger.Flush()
        }
    }
    mensaje := "Demasiados intentos fallidos desde: " + conn.RemoteAddr().String()
    helpers.WriteLog(mensaje)
}