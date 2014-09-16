package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"

	"github.com/takama/daemon"
)

const (

	// name of the service, match with executable file name
	name        = "cncore.com"
	description = "hilerchyn@gmail.com"

	// port which daemon should be listen
	port = ":9977"
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

var MonitorPanic = false

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go monitor()

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {

		if MonitorPanic {
			go monitor()
		}

		select {
		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

func monitor(){
	timer_heartbeat := time.NewTicker(5 * time.Second)

	defer func() {     //declear defer to capture panic expetion
		//fmt.Println("c")
		if err := recover(); err != nil {
			log.Println("[NetworkMonitoer]network monitor panic:",err)    //show panic error
			MonitorPanic = true
		}

		timer_heartbeat.Stop()
		//fmt.Println("d")
	}()

	//count := 1
	//Timer_Channel <- timer_heartbeat
	for true{

		select {
			//case <- time.After(time.Second*5):
			//	println("read channel timeout")
		case <-timer_heartbeat.C:
			conn, err := net.Dial("tcp", "smart.haierlife.cn:80")

			log.Println("contact: ChenTao; Tel: 15954803856; Mail: hilerchyn@gmail.com")

			if err != nil {
				// handle error

				log.Println("[NetworkMonitor] connect to failed:", err)

				if _, err = os.Stat("/opt/lampp/lampp"); err != nil {

				} else {
					cmd := exec.Command("/opt/lampp/lampp", "restart")
					err := cmd.Start()
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Waiting for command to finish...")
					err = cmd.Wait()
					log.Printf("Command finished with error: %v", err)
				}


				//continue
			} else {
				log.Println("[NetworkMonitor] success")
				conn.Close()
			}


		}

	}
}

func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)

}
