//Al finalizar reasignar la direccion del comando bash a la ubicación
//(/etc/proyOper/data/)

package reports

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func overwriteUsage() {
	updateUsage := exec.Command("bash", "/etc/proyOper/bashcommands/systemreport.sh")
	updateUsage.Stdout = os.Stdout
	updateUsage.Run()
}

func SaveReport() string {
	overwriteUsage()
	data, _ := os.ReadFile("/var/systemreport/systemreport.rpt")
	//Comando bash configurado para alojar unicamente 1 linea, por lo tanto siempre 
	//sera la 0 la que tiene el reporte
	reports := strings.Split(string(data), "\n")
	//fmt.Println("Reporte: ", reports[0])
    return reports[0]
}

func TopProcessReport() string {
	shell, err := exec.Command("bash", "/etc/proyOper/bashcommands/top_proc_by_cpu_usage.sh").CombinedOutput()
	output := string(shell)
	if err != nil {
		fmt.Println(err)
		return ""
	}
    return output
}