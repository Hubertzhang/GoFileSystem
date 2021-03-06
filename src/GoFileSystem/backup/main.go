package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "5000"
	CONN_TYPE = "tcp"
)

var sig_end chan int = make(chan int, 2)

func main() {
	err := init_config()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go TcpServer()
	go HttpServer()
	<-sig_end
	<-sig_end
}

func TcpServer() {
	defer sigEnd()

	l, err := net.Listen(CONN_TYPE, conf.Backup_ip+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Main::LISTEN::Error: " + err.Error())
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Main::ACCEPT::Error: " + err.Error())
			os.Exit(1)
		}

		go HandleConnection(conn)
	}
}

func HttpServer() {
	defer sigEnd()

	r := mux.NewRouter()
	r.HandleFunc("/kvman/countkey", HandleCountKey).Methods("GET")
	r.HandleFunc("/kvman/dump", HandleDump).Methods("GET")
	r.HandleFunc("/kvman/shutdown", HandleShutdown).Methods("GET")
	http.ListenAndServe(conf.Backup_ip+":"+conf.Http_port, r)
}

func sigEnd() {
	sig_end <- 1
}
