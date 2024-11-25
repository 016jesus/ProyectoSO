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



func GetSystemReports(intervalo time.Duration) {
		for {
	
			disk := usoDisk()
			ram := usoRAM()
			cpu := usoCPU()
	
			reporte := fmt.Sprintf("Reporte de sistema:\nDisco:\n%sMemoria:\n%sCPU:%s\n",disk, ram, cpu)
			control.Lock()
			reporteActual = reporte
			control.Unlock()
	
			// Esperar 5 segundos
			time.Sleep(intervalo)
		}
	}

	func usoCPU() string {
		cmdTop := exec.Command("top", "-b", "-n", "1")
	
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
			us := fields[3] + "%"
			return fmt.Sprint("en uso:", us, )
		} else {
			return "Error: La salida de top no tiene la estructura esperada."
		}
	}


func usoDisk() string {
	cmd := exec.Command("df", "-h", "/dev/sda5")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error obteniendo uso de disco: %s", err.Error())
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "Error: La salida de df no tiene la estructura esperada."
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return "Error: La salida de df no tiene la estructura esperada."
	}

	size, used, avail, usePercent := fields[1], fields[2], fields[3], fields[4]

	return fmt.Sprintf("/dev/sda5:\nTamaño: %s\nUsado: %s\nDisponible: %s\nPorcentaje de uso: %s\n", size, used, avail, usePercent)
}

func usoRAM() string {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error obteniendo uso de memoria: %s", err.Error())
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "Error: La salida de free no tiene la estructura esperada."
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 7 {
		return "Error: La salida de free no tiñene la estructura esperada."
	}

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