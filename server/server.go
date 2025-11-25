/*
Modulo encargado de buscar y configurar el puerto de escucha del socket, ademas de verificar
si las conexiones entrantes son validas o no
*/
package server

import (
	"os"
	"os/exec"
	"strings"
)

//Su objetivo es buscar el Listener en el archivo de configuracion para poder arrancar el socket
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

//Su objetivo es buscar si la IP que intenta conectarse si esta autorizada
func FindConnConf(fileName string, remoteIP string) bool {
	data, _ := os.ReadFile(fileName)
	lines := strings.Split(string(data), "\n")
	var allowed_ips []string
	fillData(lines, &allowed_ips)
	for _, ip := range allowed_ips {
		if ip == remoteIP {
			return true
		}
	}
	return false
}

//Llena los datos necesarios para validar
func fillData(file_lines []string, ip_list *[]string) {
	for _, conf := range file_lines {
		if conf == "" || strings.HasPrefix(conf, "#") {
			continue
		}
		if strings.HasPrefix(conf, "allowed_IP") {
			parts := strings.Split(conf, " = ")
			*ip_list = append(*ip_list, parts[1])
		}
	}
}

//Limpiar pantalla (para nada tiene que ver el nombre de la funcion con la preferencia a Windows :p)
func Cls() {
	shell := exec.Command("clear")
	shell.Stdout = os.Stdout
	shell.Run()
}
