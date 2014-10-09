package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"
	_ "bytes"

	"github.com/takama/daemon"
	"strings"

	_ "strconv"
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

	MonitorPanic := make(chan bool, 1)
	go monitor(MonitorPanic)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {


		select {
		case <-MonitorPanic:
			go monitor(MonitorPanic)
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

func monitor(MonitorPanic chan bool){
	timer_heartbeat := time.NewTicker(2 * time.Second)
	var cmd_biz *exec.Cmd
	var cmd_order *exec.Cmd

	var pid_biz int
	var pid_order int

	defer func() {     //declear defer to capture panic expetion
		//fmt.Println("c")
		if err := recover(); err != nil {
			log.Println("[NetworkMonitoer]network monitor panic:",err)    //show panic error
			MonitorPanic <- true
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

			log.Println("contact: ChenTao; Tel: 15954803856; Mail: hilerchyn@gmail.com")
			//log.Println(os.Getwd())
			os.Chdir("/root/printer_upload")
			//log.Println(os.Getwd())

			look, _:=exec.LookPath("/bin/")
			log.Println(look)

			cmd := exec.Command("ps", "aux")
			out, _ := cmd.Output()
			result := string(out);
			//log.Println(result)

			count_biz := strings.Count(result, "biz-server")
			log.Println("count:",count_biz)
			if count_biz <=0 {
				//look, _:=exec.LookPath("/usr/bin/")
				//log.Println(look)
				cmd_biz= exec.Command("/usr/bin/nohup","/root/printer_upload/biz-server")
				cmd_biz.Start()
				log.Println(cmd_biz.Process.Pid)
				pid_biz = cmd_biz.Process.Pid
			}

			count_biz_kill := strings.Count(result, "[biz-server]")
			log.Println("count kill:",count_biz_kill)
			if count_biz_kill >0 {
				p, err := os.FindProcess(pid_biz)
				p.Wait()
				log.Println(err)
			}

			count_order := strings.Count(result, "order-server")
			log.Println("count:",count_order)
			if count_order <=0 {
				//look, _:=exec.LookPath("/usr/bin/")
				//log.Println(look)
				cmd_order = exec.Command("/usr/bin/nohup","/root/printer_upload/order-server")
				cmd_order.Start()
				log.Println(cmd_order.Process.Pid)
				pid_order = cmd_order.Process.Pid
				//log.Println(err)
			}

			count_order_kill := strings.Count(result, "[order-server]")
			log.Println("count kill:",count_order_kill)
			if count_order_kill >0 {
				p, err := os.FindProcess(pid_order)
				p.Wait()
				log.Println(err)
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
