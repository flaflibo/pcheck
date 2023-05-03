package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var value = "#99@123"
var erroMsg []string
var portMap map[int]bool
var wg sync.WaitGroup

func startUdpServer(port int) {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Not able to create the udp server")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Not able to create the udp server", err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	for {

		_, _, err = conn.ReadFromUDP(buf)
		if err != nil {
			// log.Fatal("could not read data")
			continue
		}

		portMap[addr.Port] = true
		wg.Done()
		break

	}

}

func printStats() {
	log.Println("Could not receive anything on these ports: ")
	for key, value := range portMap {
		if !value {
			log.Println(key)
		}
	}
}

func sendUdpMessage(ip string, port int) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal("could not create udp addr")
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Println("could not connect to the udp server")
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(value))
	if err != nil {
		log.Fatal("could not sendsend the message")
	}

}

func main() {
	start := flag.Int("start", 0, "start port")
	end := flag.Int("end", 0, "end port")
	tcp := flag.Bool("t", false, "if true use tcp otherwise udp")
	isServer := flag.Bool("s", false, "if true is server otherwise client")
	ip := flag.String("ip", "", "server ip")
	flag.Parse()

	portMap = make(map[int]bool)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	if *isServer {

		k := 0

		for i := *start; i < *end; i++ {
			log.Println("start port --> ", i)
			if *tcp {

			} else {
				wg.Add(1)
				go startUdpServer(i)
				k++
				if k == 10 {
					//after we started 10 servers we will wait for them to finish
					//so we can start the next set of servers
					wg.Wait()
					k = 0
				}
			}
		}

		wg.Wait()
		printStats()

		// <-sigChannel

	} else {
		//is client

		for {
			select {
			case <-sigChannel:
				return
			default:
				for i := *start; i < *end; i++ {
					if *tcp {

					} else {
						log.Println("Going to send to port ", i)
						sendUdpMessage(*ip, i)
					}
				}
			}

		}
	}

}
