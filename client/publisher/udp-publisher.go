package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"time"
)

const maxMessageSize int = 65507

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer conn.Close()

	dummyData := make([]byte, 100000)
	for i := range dummyData {
		dummyData[i] = byte(i)
	}

	headerSize := 1 + 2 + 4*3 // header 14bytes
	contentSize := maxMessageSize - headerSize

	message := make([]byte, maxMessageSize)

	messageIndex := 0

	for {
		//message := []byte{1, 1, 2, 3}
		messageIndex++
		data := dummyData
		packetNum := (len(data) + contentSize - 1) / contentSize
		for i := 0; i < packetNum; i++ {
			messageSize := headerSize
			if i < packetNum-1 {
				messageSize += contentSize
			} else {
				messageSize += len(data) % contentSize
			}
			message[0] = 2 //publish

			buf := new(bytes.Buffer)
			// メッセージサイズ
			binary.Write(buf, binary.LittleEndian, int16(messageSize))
			// index
			err := binary.Write(buf, binary.LittleEndian, int32(messageIndex))
			// パケットインデックス
			binary.Write(buf, binary.LittleEndian, int32(i))
			// パケット数
			binary.Write(buf, binary.LittleEndian, int32(packetNum))
			//copy(message[3:7], buf.Bytes())
			copy(message[1:], buf.Bytes())
			copy(message[headerSize:], data[i*contentSize:i*contentSize+messageSize-headerSize])

			//_, err := conn.Write(dummy_data)
			_, err = conn.Write(message)

			//_, err := conn.Write(message)
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
