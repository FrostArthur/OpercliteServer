/*Modulo encargado de obtener información del sistema para enviarla al socket cliente*/
package reports

import (
	"fmt"
	"os/exec"
)

//Ejecuta la funcion bash para sobreescribir el uso del system
func OverwriteUsage() string {
	updatedUsage, err := exec.Command("bash", "/etc/proyOper/bashcommands/systemreport.sh").CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	output := string(updatedUsage)
	return output
}

func TopProcessReport() string {
	shell, err := exec.Command("bash", "/etc/proyOper/bashcommands/top_proc_by_cpu_usage.sh").CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	output := string(shell)
    return output
}