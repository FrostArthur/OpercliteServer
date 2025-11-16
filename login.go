//Al finalizar reasignar la lectura del vector de strings "data" a donde se ubica principalmente el
//archivo proyOper.users -> (/etc/proyOper/data/)

package login

import (
	"os"
	"strings"
	"fmt"
)

func FindCli(user string, passwd string) bool {
	data, _ := os.ReadFile("/etc/proyOper/data/proyOper.users")
	credentials := strings.Split(string(data), "\n")
	var usercli []string
	var hpasswdcli []string
	fillData(credentials, &usercli, &hpasswdcli)
    for x, cli := range usercli {
		if (user == cli) {
			if (passwd == hpasswdcli[x]) {
				return true
			}
		}
	}
	return false
}

func fillData(users []string, usercli *[]string, hpasswdcli *[]string) {
	for _, cli := range users {
		cli = strings.TrimSpace(cli)
		if cli == " " || strings.HasPrefix(cli, "#") {
			continue
		}
		parts := strings.Split(cli, ":") 
		//Para evitar errores por usar direcciones de memoria invalidas
		if (len(parts) < 2) {
			continue
		}
		fmt.Println(len(parts), "cadenas")
		*usercli = append(*usercli, parts[0])
		*hpasswdcli = append(*hpasswdcli, parts[1])
	}
}
