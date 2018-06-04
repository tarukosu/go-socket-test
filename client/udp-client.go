package main

import (
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		recvBuf := make([]byte, 65507)
		for {
			n, err := conn.Read(recvBuf)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}

			log.Printf("Received data: %s", string(recvBuf[:n]))
		}
	}()

	for {
		//message := make([]byte, 65507)
		message := make([]byte, 65507)

		//n, err := conn.Write([]byte("Ping"))
		_, err := conn.Write(message)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		/*
			if len(buf) != n {
				log.Printf("data size is %d, but sent data size is %d", len(buf), n)
			}
		*/

		/*
			for {
				recvBuf := make([]byte, 65507)
				n, err = conn.Read(recvBuf)
				if err != nil {
					log.Fatalln(err)
					os.Exit(1)
				}

				log.Printf("Received data: %s", string(recvBuf[:n]))
			}
		*/
		time.Sleep(1 * time.Second)
	}
}
