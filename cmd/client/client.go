package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	CONN_PORT = ":3333"
	CONN_TYPE = "tcp"

	MSG_DISCONNECT = "Disconnected from the server.\n"
)

var (
	wg       sync.WaitGroup
	addrFlag = flag.String("addr", "", "Show debug information.")
	NameFlag = flag.String("name", "", "Path to the log file.")
)

// Reads from the socket and outputs to the console.
func Read(conn net.Conn) {
	flag.Parse()
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(MSG_DISCONNECT)
			wg.Done()
			return
		}
		if str == "> There is a user which is used for this name. Please choose another name\n" {
			fmt.Println("There is a user in server which is using same name with you.\nPlease choose another name.")
			os.Exit(1)
		}
		fmt.Print(str)
	}
}

// Reads from Stdin, and outputs to the socket.
func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = writer.WriteString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
func sendName(conn net.Conn, name string) error {
	writer := bufio.NewWriter(conn)

	str := "/name " + name + "\n"

	_, err := writer.WriteString(str)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func quit(conn net.Conn) error {
	writer := bufio.NewWriter(conn)

	str := "/quit\n"

	_, err := writer.WriteString(str)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Starts up a read and write thread which connect to the server through the
// a socket connection.
func main() {
	flag.Parse()
	wg.Add(2)
	if *NameFlag == "" {
		fmt.Println("Pls write a name\nExample:\n\t-name Selahattin")
		os.Exit(1)
	}
	if *addrFlag == "" {
		fmt.Println("Pls write a server adress\nExample:\n\t-addr localhost:8080\n\t-addr :8080\n\t-addr 127.0.0.1:8080")
		os.Exit(1)
	}
	conn, err := net.Dial(CONN_TYPE, *addrFlag)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			quit(conn)
			os.Exit(1)
		}()
	}()
	go Read(conn)
	go Write(conn)
	sendName(conn, *NameFlag)
	wg.Wait()

}
