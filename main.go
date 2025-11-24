package main
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"OpercliteServer/login"
	"OpercliteServer/reports"
	"OpercliteServer/server"
	"time"
)

func main() {
	//Inicializacion del servidor
	const dbase string = "/etc/operserver.conf"
	server.Cls()
	fmt.Println("\n**************************************")
	fmt.Println("*           Operclite Server           *")
	fmt.Println("****************************************")
	localPort := server.StartServer(dbase)
	dirTCP, _ := net.ResolveTCPAddr("tcp4", ":"+localPort)
	fmt.Println("Abriendo el puerto ", localPort)
	socketServer, _ := net.ListenTCP("tcp4", dirTCP)
	fmt.Println("Esperando conexiones...")
	//Socket remoto conectandose
	remoteSocket, _ := socketServer.Accept()
	cli_IP, _, _ := net.SplitHostPort(remoteSocket.RemoteAddr().String())
	fmt.Println("Solicitud de: ", cli_IP)
	//Validacion de IP
	if !server.FindConnConf(dbase, cli_IP) {
		sendDefaultMsg(&remoteSocket, "Unauthorized")
		server.Cls()
		fmt.Println("Unauthorized, cerrando conexion")
		remoteSocket.Close()
		return
	}
	sendDefaultMsg(&remoteSocket, "Sucessful")
	fmt.Println("OK!")
	str_interval := recCredentialMsg(&remoteSocket)
	interval, err := strconv.Atoi(str_interval)
	if err != nil {
		fmt.Println("Error al convertir: ", err)
		return
	}
	//Login
	var attempts int = 0
	var counter int = 1
	//Le enviamos un mensaje que diga freeze: <segundos>, en caso de que haya mas de 10 intentos, cerrar la conexion
	for {
		fmt.Println("Pidiendo usuario...")
		sendDefaultMsg(&remoteSocket, "Ingrese su usuario a continuacion")
		userclient := recCredentialMsg(&remoteSocket) //Espera recibir el usuario
		fmt.Println("Pidiendo contraseña...")
		sendDefaultMsg(&remoteSocket, "Ingrese su contraseña a continuacion")
		passwdclient := recCredentialMsg(&remoteSocket) //Espera recibir contraseña
		if login.CliFound(userclient, passwdclient) {
			sendDefaultMsg(&remoteSocket, "Login OK")
			break
		}
		attempts++
		if attempts == 6 {
            sendDefaultMsg(&remoteSocket, "Locked")
			server.Cls()
			remoteSocket.Close()
			fmt.Println("Conexion finalizada, error de login")
			return
		}
		if attempts >= 3 { 
            counter += counter
			fmt.Println("Muchos intentos por parte del cliente, poner en espera durante ", counter)
			sendDefaultMsg(&remoteSocket, "Freeze: " + strconv.Itoa(counter))
			continue
		}
		sendDefaultMsg(&remoteSocket, "Retry")
	}
	attempts = 0
	counter = 0
	fmt.Println("Login OK")
	//Comunicacion principal
	quitServer := make(chan bool)
	go recMsg(remoteSocket, cli_IP, quitServer)
	go sendUsageReport(remoteSocket, quitServer, interval)
	<-quitServer
	server.Cls()
	fmt.Println("Gracias por usar el servidor de Operclite :)")
	remoteSocket.Close()
	os.Exit(0)
}

func recMsg(socket net.Conn, remoteIP string, quitFlag chan bool) {
	for {
		msg, _ := bufio.NewReader(socket).ReadString('\n')
		fmt.Println("Client [ ", remoteIP, " ] --> ", msg)
		if strings.HasPrefix(msg, "exec:") {
			fmt.Println("Tratamos con un comando")
			command := strings.TrimSpace(strings.TrimPrefix(msg, "exec:"))
			fmt.Println(command)
			execCommand(socket, command)
			continue
		}
		if msg == "bye\n" { //Formato \r\n para recibir correctamente del módulo de Windows
			quitFlag <- true
			return
		}
	}
}

func sendUsageReport(socket net.Conn, quitFlag chan bool, interval int) {
	numReport := 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		numReport++
		msg := "Report: " + reports.OverwriteUsage() + "\r\n" //Formato \r\n para recibir correctamente del módulo de Windows
		env := bufio.NewWriter(socket)
		env.WriteString(msg)
		env.Flush()
		socket.Write([]byte(reports.TopProcessReport()))
		if msg == "bye\n" {
			quitFlag <- true
			return
		}
	}
}

func sendDefaultMsg(socket *net.Conn, defMsg string) {
	env := bufio.NewWriter(*socket)
	env.WriteString(defMsg + "\r\n") //Formato \r\n para recibir correctamente del módulo de Windows
	env.Flush()
}

func recCredentialMsg(socket *net.Conn) string {
	msg, _ := bufio.NewReader(*socket).ReadString('\n')
	msg = strings.TrimSpace(msg)
	fmt.Println("Client --> ", msg)
	var bInterval bool
	msg, bInterval = strings.CutPrefix(msg ,"Time: ")
	if bInterval {
		fmt.Println("Intervalo de tiempo recibido :", msg)
	}
	return msg
}

func execCommand(socket net.Conn, commandBash string) {
	shell, err := exec.Command("bash", "-c", commandBash).CombinedOutput()
	output := string(shell)
	if err != nil {
		output += fmt.Sprintf("\nError al ejecutar: %v\n", err)
	}
	lines := strings.Count(output, "\n")

	if lines > 50 {
		socket.Write([]byte(fmt.Sprintf("El comando genera %d líneas. No se mostrará completo.\n", lines)))
		parts := strings.SplitN(output, "\n", 41)
		output = strings.Join(parts[:40], "\n") + "\n... [Salida truncada]\n"
	} else {
		fmt.Println("Si estoy apto para enviarte comandos mi rey :3")
	}
    
	socket.Write([]byte(output + "\r\n"))
}
