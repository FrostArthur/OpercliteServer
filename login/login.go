/*Modulo encargado de validar credenciales para permitir o no el acceso a la comunicacion principal
entre los dos sockets*/
package login

import (
	"os"
	"strings"
)

//Valida contra la base de datos de usuarios del socket
func CliFound(user string, passwd string) bool {
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

//Llena los datos necesarios para validar
func fillData(users []string, usercli *[]string, hpasswdcli *[]string) {
	for _, cli := range users {
		cli = strings.TrimSpace(cli)
		//Evita lineas invalidas
		if cli == " " || strings.HasPrefix(cli, "#") {
			continue
		}
		parts := strings.Split(cli, ":") 
		//Para evitar errores por usar direcciones de memoria invalidas
		if (len(parts) != 2) {
			continue
		}
		*usercli = append(*usercli, parts[0])
		*hpasswdcli = append(*hpasswdcli, parts[1])
	}
}
