package helpers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Ram string
var Cpu string
var Disk string



func GetSystemReports(intervalo time.Duration) {
		for {
	
			Disk = usoDisk()
			Ram = usoRAM()
			Cpu = usoCPU()
	
			reporte := fmt.Sprintf("Reporte de sistema:\nDisco:\n%sMemoria:\n%sCPU:%s\n",Disk, Ram, Cpu)
			control.Lock()
			reporteActual = reporte
			control.Unlock()
	
			// Esperar 5 segundos
			time.Sleep(intervalo)
		}
	}

	func usoCPU() string {
		cmdTop := exec.Command("/bin/sh", "-c", "top -bn1")	
		var outTop bytes.Buffer
		cmdTop.Stdout = &outTop
		err := cmdTop.Run()
		if err != nil {
			return fmt.Sprint("Error ejecutando top:", err)
		}
	
		// Filtramos la salida con `grep`
		cmdGrep := exec.Command("grep", "%Cpu")
	
		// Pasamos la salida de `top` al comando `grep`
		cmdGrep.Stdin = &outTop
		var outGrep bytes.Buffer
		cmdGrep.Stdout = &outGrep
		err = cmdGrep.Run()
		if err != nil {
			return fmt.Sprint("Error ejecutando grep:", err)
		}
	
		// Usamos strings para dividir y manipular la salida
		output := outGrep.String()
		fields := strings.Fields(output)
	
		// Formateamos la salida
		if len(fields) >= 8 {
			us := fields[3]
			us = strings.Trim(us, "\n")
			us = strings.Replace(us, ",", ".", 1)
			return "\nen uso:" + us + "%" + "\n"
		} else {
			return "Error: La salida de top no tiene la estructura esperada."
		}
	}


func usoDisk() string {

	cmd := exec.Command("df", "-h", "/dev/sda5")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	fields := strings.Fields(lines[1])
	size, used, avail, usePercent := fields[1], fields[2], fields[3], fields[4]

	return fmt.Sprintf("/dev/sda5:\nTamaño: %s\nUsado: %s\nDisponible: %s\nPorcentaje de uso: %s\n", size, used, avail, usePercent)
}

func usoRAM() string {
	cmd := exec.Command("free", "-m")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	fields := strings.Fields(lines[1])

	t, u, f := fields[1], fields[2], fields[3]
	total, _ := strconv.ParseFloat(t, 64)
	used, _ := strconv.ParseFloat(u, 64)
	free, _ := strconv.ParseFloat(f, 64)
	usedPercent := (used / total) * 100

	return fmt.Sprintf("Total: %.2f MB \nUsada: %.2f MB \nLibre: %.2f MB \nPorcentaje de uso: %.2f%%\n", total, used, free, usedPercent)
}


func getOutput(remoteReader *bufio.Reader) strings.Builder {
	var output strings.Builder
	for {
		linea, err := remoteReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			output.WriteString(fmt.Sprintf("Error leyendo la salida remota: %v\n", err))
			break;
		}
		if strings.TrimSpace(linea) == "" {
			break
		}
		output.WriteString(linea)
	}
	return output
}