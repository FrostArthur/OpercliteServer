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
	const dbase string = "/etc/operserver.conf"
	server.Cls()
	fmt.Println("\n**************************************")
	fmt.Println("*           Operclite Server           *")
	fmt.Println("****************************************")
	localPort := server.StartServer(dbase)
	remoteSocket, err := connectionPhase(localPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cli_IP, _, _ := net.SplitHostPort(remoteSocket.RemoteAddr().String())
	fmt.Println("Solicitud de: ", cli_IP)
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
	var attempts int = 0
	var counter int = 1
	loginProcess(remoteSocket, attempts, counter)
	attempts = 0
	counter = 0
	fmt.Println("Login OK")
	quitServer := make(chan bool)
	go recMsg(remoteSocket, cli_IP, quitServer)
	go sendUsageReport(remoteSocket, quitServer, interval)
	<-quitServer
	server.Cls()
	fmt.Println("Gracias por usar el servidor de Operclite :)")
	remoteSocket.Close()
	os.Exit(0)
}

func connectionPhase(localPort string) (net.Conn, error) {
	dirTCP, errAddr := net.ResolveTCPAddr("tcp4", ":" + localPort)
	if errAddr != nil {
		return nil, errAddr
	}
	fmt.Println("Abriendo el puerto ", localPort)
	serverSocket, errListen := net.ListenTCP("tcp4", dirTCP)
	if errListen != nil {
		return nil, errListen
	}
	fmt.Println("Esperando conexiones...")
	remoteSocket, errAccept := serverSocket.Accept()
	if errAccept != nil {
		return nil, errAccept
	}
	return remoteSocket, nil
}

func loginProcess(remoteSocket net.Conn, attempts int, counter int) {
	for {
		fmt.Println("Pidiendo usuario...")
		sendDefaultMsg(&remoteSocket, "Ingrese su usuario a continuacion")
		userclient := recCredentialMsg(&remoteSocket) 
		fmt.Println("Pidiendo contraseña...")
		sendDefaultMsg(&remoteSocket, "Ingrese su contraseña a continuacion")
		passwdclient := recCredentialMsg(&remoteSocket) 
		client := strings.TrimSpace(userclient + ":" + passwdclient)
		if login.CliFound(client) {
			sendDefaultMsg(&remoteSocket, "Login OK")
			return
		}
		attempts++
		if attempts == 6 {
            sendDefaultMsg(&remoteSocket, "Locked")
			server.Cls()
			remoteSocket.Close()
			fmt.Println("Conexion finalizada, error de login")
			os.Exit(0)
		}
		if attempts >= 3 { 
            counter += counter
			fmt.Println("Muchos intentos por parte del cliente, poner en espera durante ", counter)
			sendDefaultMsg(&remoteSocket, "Freeze: " + strconv.Itoa(counter))
			continue
		}
		sendDefaultMsg(&remoteSocket, "Retry")
	}
}

func recMsg(socket net.Conn, remoteIP string, quitFlag chan bool) {
	for {
		socket.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		msg, err := bufio.NewReader(socket).ReadString('\n')
		if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Timeout() {
			    continue
		    }
			fmt.Println("Error de conexion: ", err)
			quitFlag <- true
			os.Exit(1)
		}
		fmt.Println("Client [ ", remoteIP, " ] --> ", msg)
		if strings.HasPrefix(msg, "exec:") {
			fmt.Println("Tratamos con un comando")
			env := bufio.NewWriter(socket)
			command := strings.TrimSpace(strings.TrimPrefix(msg, "exec:"))
			fmt.Println(command)
			env.WriteString(execCommand(command))
			env.Flush()
			continue
		}
		if msg == "bye\n" { 
			quitFlag <- true
			return
		}
	}
}

func sendUsageReport(socket net.Conn, quitFlag chan bool, interval int) {
	env := bufio.NewWriter(socket)
	numReport := 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		numReport++
		msg := "Report: " + reports.OverwriteUsage() + "\r\n"
		env.WriteString(msg)
		env.Flush()
		env.WriteString(reports.TopProcessReport())
		env.Flush()
		if msg == "bye\n" {
			quitFlag <- true
			return
		}
	}
}

func sendDefaultMsg(socket *net.Conn, defMsg string) {
	env := bufio.NewWriter(*socket)
	env.WriteString(defMsg + "\r\n")
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

func execCommand(commandBash string) string {
	shell, err := exec.Command("bash", "-c", commandBash).CombinedOutput()
	output := string(shell)
	if err != nil {
		output += fmt.Sprintf("\nError al ejecutar: %v\n", err)
	}
	lines := strings.Count(output, "\n")

	if lines > 50 {
		parts := strings.SplitN(output, "\n", 51)
		output = strings.Join(parts[:50], "\n") + "\n... [Salida recortada]\n"
	}
	return output
}
