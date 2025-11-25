package login

import (
	"os"
	"strings"
)

func CliFound(client string) bool {
	data, _ := os.ReadFile("/etc/proyOper/data/proyOper.users")
	lines := strings.Split(string(data), "\n")
	var usercli []string
	var hpasswdcli []string
	fillData(lines, &usercli, &hpasswdcli)
    for x, user := range usercli {
		if (strings.HasPrefix(client, user)) {
			return client == strings.TrimSpace(user + ":" + hpasswdcli[x])
		}
	}
	return false
}

func fillData(file_lines []string, usercli *[]string, hpasswdcli *[]string) {
	for _, cli := range file_lines {
		cli = strings.TrimSpace(cli)
		if cli == " " || strings.HasPrefix(cli, "#") {
			continue
		}
		parts := strings.Split(cli, ":") 
		if (len(parts) != 2) {
			continue
		}
		*usercli = append(*usercli, parts[0])
		*hpasswdcli = append(*hpasswdcli, parts[1])
	}
}

