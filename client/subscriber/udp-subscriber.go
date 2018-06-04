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

			log.Printf("Data Received")
			log.Println(n)
			for i := 0; i < 10; i++ {
				log.Printf("%x", recvBuf[i])
			}

			//log.Printf("Received data: %s", string(recvBuf[:n]))
			//log.Println(n)
			//log.Printf("Received data: %s", string(recvBuf[:10]))
		}
	}()

	for {
		//message := []byte("subscribe")
		message := []byte{2}

		_, err := conn.Write(message)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		time.Sleep(1 * time.Second)
	}
}
