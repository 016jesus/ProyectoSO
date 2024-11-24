package helpers
import ("bytes"
	"fmt"
	"os/exec"
	"strings"
	"bufio"
	"time"
	"io"
)



func GetSystemReports(intervalo time.Duration) {
		for {
	
			usoDisco := usoDisk()
			usoMemoria := usoRAM()
			usoCPU := usoCPU()
	
			reporte := fmt.Sprintf(
				"Reporte de sistema:\nDisco: %s\nMemoria: %s\nCPU: %s\n",
				usoDisco, usoMemoria, usoCPU,
			)
	
			// Actualizar el reporte compartido
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

	// Ahora filtramos la salida con `grep`
	cmdGrep := exec.Command("grep", "Cpu(s)")

	// Pasamos la salida de `top` al comando `grep`
	cmdGrep.Stdin = &outTop
	var outGrep bytes.Buffer
	cmdGrep.Stdout = &outGrep
	err = cmdGrep.Run()
	if err != nil {
		return fmt.Sprint("Error ejecutando grep:", err)
	}

	// Filtramos la salida con `awk`
	// Usamos strings para dividir y manipular la salida
	output := outGrep.String()
	fields := strings.Fields(output)

	// Formateamos la salida
	if len(fields) >= 8 {
		us := fields[1]
		sy := fields[3]
		ni := fields[5]
		id := fields[7]

		// Imprimimos el resultado con el formato que queremos
		return fmt.Sprintf("us: %s sy: %s ni: %s id: %s\n", us, sy, ni, id)
	} else {
		return "Error: La salida de top no tiene la estructura esperada."
	}
}


func usoDisk() string {
cmd := exec.Command("df", "-h")
output, err := cmd.Output()
if err != nil {
	return fmt.Sprintf("Error obteniendo uso de disco: %s", err.Error())
}
return strings.TrimSpace(string(output))
}

func usoRAM() string {
cmd := exec.Command("free", "-m")
output, err := cmd.Output()
if err != nil {
	return fmt.Sprintf("Error obteniendo uso de memoria: %s", err.Error())
}
return strings.TrimSpace(string(output))
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