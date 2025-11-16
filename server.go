//Manipular despues para mejorar seguridad ofreciendo un timeoff al cliente
//Esto esta especificado en ponerle una fecha de vencimiento del timeoff cada vez
//que el cliente en blacklist quiera ingresar, realizar esas configuraciones en este
//archivo

package server

import (
	//"fmt"
	"os"
	"os/exec"
	"strings"
	//"time"
	//"sync"
)

func StartServer(file string) string {
	data, _ := os.ReadFile(file)
	confText_line := strings.Split(string(data), "\n")
	var listener string
	for _, line := range confText_line {
		if strings.HasPrefix(line, "listener") {
			listenerpart := strings.Split(line, " = ")
			listener = listenerpart[1]
			break
		}
	}
	return listener
}

func FindConnConf(fileName string, remoteIP string) bool {
	data, _ := os.ReadFile(fileName)
	credentials := strings.Split(string(data), "\n")
	var allowed_ips []string
	fillData(credentials, &allowed_ips)
	for _, ip := range allowed_ips {
		if ip == remoteIP {
			return true
		}
	}
	return false
}

func fillData(conf_list []string, ip_list *[]string) {
	for _, conf := range conf_list {
		if conf == "" || strings.HasPrefix(conf, "#") {
			continue
		}
		if strings.HasPrefix(conf, "allowed_IP") {
			parts := strings.Split(conf, " = ")
			*ip_list = append(*ip_list, parts[1])
		}
	}
}

func Cls() {
	shell := exec.Command("clear")
	shell.Stdout = os.Stdout
	shell.Run()
}
