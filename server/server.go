package server

import (
	"os"
	"os/exec"
	"strings"
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

func Cls() {
	shell := exec.Command("clear")
	shell.Stdout = os.Stdout
	shell.Run()
}
