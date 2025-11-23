/*Modulo encargado de registrar logs para controlar la cantidad de intentos que llevan los clientes
y manejando persistencia en caso de que llegue a cerrar abruptamente el programa*/
package logs

import (
	"os"
	"strings"
	"fmt"
	"strconv"
)

/*Esta es la estructura que maneja los datos del cliente (son por IP, no por usuarios) que se esta intentando 
logear actualmente*/
type TimeOffData struct {
	ip *string
    timeoff *int64
	attempts *int64
}

//Es el respaldo de los datos del archivo plano logs.log que garantiza que una vez se cierre el programa se sobreescriban toda la información
type BackupData struct {
	clients_ip *[]string
    clients_timeoff *[]string
	clients_attempts *[]string
}

func GetAttempts(data TimeOffData) int64 {
    return *data.attempts
}

func GetTimeOff(data TimeOffData) int64 {
    return *data.timeoff
}

//Carga la información de tiempo de espera que tiene un cliente si el programa cerró previamente durante un tiempo de espera
func LoadTimeOffData(client_ip string) (TimeOffData, BackupData) {
	timeData := TimeOffData {}
    backup := BackupData {}
	fileData_B, err := os.ReadFile("../logs.log")
	if err != nil {
		fmt.Println("Error al abrir el archivo: ", err)
        return timeData, backup
	}
    fileData := strings.Split(string(fileData_B), "\n")
	var timeOff, attempts int64
	var ipsList, timeoffList, attemptsList []string
	errInFill := fillData(fileData, client_ip, &timeOff, &attempts, &ipsList, &timeoffList, &attemptsList)
	if errInFill != nil {
       fmt.Println("Error en el llenado de datos: ", err)
	   return timeData, backup
	}
	timeData = TimeOffData {
		ip: &client_ip,
		timeoff: &timeOff,
		attempts: &attempts,
	}
	backup = BackupData {
        clients_ip: &ipsList,
		clients_timeoff: &timeoffList,
		clients_attempts: &attemptsList,
	}
	return timeData, backup
}

func fillData(fileData []string, client_ip string, 
	timeOff *int64, attempts *int64, ipList *[]string, 
	timeOffList *[]string, attemptList *[]string) error {
	for _, client := range fileData {
		client = strings.TrimSpace(client)
        //Se ignora la linea que este vacia o comentada
		if client == "" || strings.HasPrefix(client, "#") {
			continue
		}
		data := strings.Split(client, ":")
		//Comprobamos si estan los 3 campos que se necesitan
		if (len(data) != 3) {
			fmt.Println("Error, linea corrupta")
			continue
		}
		//Cargamos el respaldo de los datos del archivo logs.log
        *ipList = append(*ipList, data[0])
		*timeOffList = append(*timeOffList, data[1])
		*attemptList = append(*attemptList, data[2])
		/*Comprobamos si el cliente efectivamente esta en la lista de bloqueados con la
		ip que se recibe por parametro (se comprueba siempre en la primera posicion)*/
		if (data[0] == client_ip) { //Si coincide se sabe que [1] y [2] contienen el tiempo de espera y los intentos respectivamente
            tOff, err := strconv.ParseInt(data[1], 10, 64)
		    if err != nil {
			   return err
		    }
			*timeOff = tOff
			atts, err := strconv.ParseInt(data[2], 10, 64)
		    if err != nil {
			   return err
		    }
			*attempts = atts
		}
	}
	return nil
}

//Valida si alguna de las dos estructuras esta vacia
func CurrentDataIsEmpty(currentData TimeOffData) bool {
    return currentData == TimeOffData {}
}

func BackupIsEmpty(currentBackup BackupData) bool {
    return currentBackup == BackupData {}
}

func SaveTimeOffData() {

}

