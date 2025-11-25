/*Modulo encargado de validar credenciales para permitir o no el acceso a la comunicacion principal
entre los dos sockets*/
package login

import (
	"os"
	"strings"
)

//Valida contra la base de datos de usuarios del socket
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

//Llena los datos necesarios para validar
func fillData(file_lines []string, usercli *[]string, hpasswdcli *[]string) {
	for _, cli := range file_lines {
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

