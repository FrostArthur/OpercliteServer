package reports

import (
	"fmt"
	"os/exec"
)

func OverwriteUsage() string {
	updatedUsage, err := exec.Command("bash", "/etc/proyOper/bashcommands/systemreport.sh").CombinedOutput()
	output := string(updatedUsage)
	if err != nil {
		fmt.Println(err)
		return ""
	}
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