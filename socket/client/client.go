package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
    "time"
)

const (
	// 1 byte length stopCharacter
	stopCharacter       = byte(4)
	stopCharacterLength = 1
	exitCommand         = "exit"

	// 1 byte length data type
	// Exit. type:1, content:""
	// Chat. type:2, content:chat message
	typeExit = byte(1)
	typeChat = byte(2)

	// 16 byte length nickname
	nicknameLen        = 16
	NicknameNotInputed = ""
)

func randomNickname(len int) []byte {
	bytes := make([]byte, len)
    rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // Atoz in ascii
	}
	return bytes
}

func SocketClient(dest string, nickname string) {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()
	defer log.Printf("Close connection to %v", conn.RemoteAddr())

	log.Printf("Connected to %v from %v", conn.RemoteAddr(), conn.LocalAddr())

	if strings.Compare(nickname, NicknameNotInputed) == 0 || len(nickname) > 16 {
		nickname = string(randomNickname(nicknameLen))
	}

	var (
		buf  = make([]byte, 1024)
		r    = bufio.NewReader(conn)
		w    = bufio.NewWriter(conn)
		done = false
	)

	go func() {
		for {
			n, err := r.Read(buf)

			switch err {
			case io.EOF:
				return
			case nil:
				if buf[n-stopCharacterLength] == stopCharacter {
					log.Printf("Receive: %s", buf[:n-stopCharacterLength])
				} else {
					log.Printf("Receive message without stop character: %v", buf[:n])
				}
			default:
				if done {
					return
				}
				log.Fatalf("Receive data failed:%s", err)
				return
			}
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		w.Write([]byte(nickname + ": "))
		w.Write([]byte(scanner.Text()))
		w.Write([]byte{stopCharacter})
		w.Flush()
		log.Printf("Send: %v", scanner.Text())
		if strings.Compare(scanner.Text(), exitCommand) == 0 {
			log.Println("Client Exit!")
			done = true
			return
		}
	}
}

func main() {
	var (
		dest     = flag.String("d", "127.0.0.1:30000", "Server address in IP:PORT format")
		nickname = flag.String("n", NicknameNotInputed, "User nickname")
	)

	//var logfile, err = os.OpenFile("./client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	//if err != nil {
	//    log.Fatalf("error opening file: %v", err)
	//}
	//defer logfile.Close()
	//log.SetOutput(logfile)
	log.Printf("---Start logging---")

	flag.Parse()
	SocketClient(*dest, *nickname)
}
