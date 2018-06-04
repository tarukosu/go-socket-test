package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Subscribers struct {
	mu          sync.Mutex
	subscribers []Subscriber
}

func (subs *Subscribers) lookUpSubscriber(addr *net.UDPAddr) (sub *Subscriber, ok bool) {
	subs.mu.Lock()
	defer subs.mu.Unlock()
	for i, s := range subs.subscribers {
		if s.isSamePort(addr) {
			return &subs.subscribers[i], true
		}
	}
	return nil, false
}

func (subs *Subscribers) addMember(addr *net.UDPAddr) {
	subs.mu.Lock()
	defer subs.mu.Unlock()

	subscriber := Subscriber{addr: addr}
	subscriber.update()
	subs.subscribers = append(subs.subscribers, subscriber)

}

func (subs *Subscribers) removeMember(s Subscriber) {

	subs.mu.Lock()
	defer subs.mu.Unlock()
	newSubscribers := []Subscriber{}

	for _, sub := range subs.subscribers {
		if !sub.isSame(s) {
			newSubscribers = append(newSubscribers, sub)
		}
	}
	subs.subscribers = newSubscribers
}

func (subs *Subscribers) checkTimeOut(duration time.Duration) {
	log.Println("check time out")
	subs.mu.Lock()
	timeoutPeers := []Subscriber{}

	timeoutTime := time.Now().Add(-duration)

	for _, s := range subs.subscribers {
		log.Println(s.updated)
		if s.updated.Before(timeoutTime) {
			timeoutPeers = append(timeoutPeers, s)
			log.Printf("timeout: %s", s)
		}
	}

	subs.mu.Unlock()

	for _, s := range timeoutPeers {
		subs.removeMember(s)
	}
}

type Subscriber struct {
	addr    *net.UDPAddr
	updated time.Time
}

func (s Subscriber) String() string {
	return fmt.Sprintf("%s, update: %s", s.addr.String(), s.updated)
}

func (s *Subscriber) update() {
	s.updated = time.Now()
}

func (s Subscriber) isSamePort(addr *net.UDPAddr) bool {
	return s.addr.IP.String() == addr.IP.String() && s.addr.Port == addr.Port
}

func (s Subscriber) isSame(other Subscriber) bool {
	return s.addr.IP.String() == other.addr.IP.String() && s.addr.Port == other.addr.Port
}

var subscribers = Subscribers{}

//subscribers := make([]Subscriber, 0)

func checkTimeOut() {
	go func() {
		for {
			subscribers.checkTimeOut(30 * time.Second)
			time.Sleep(1 * time.Second)
		}
	}()
}

func main() {
	checkTimeOut()

	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
	updLn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	//subscribers := make([]Subscriber, 0)

	buf := make([]byte, 65507)
	log.Println("Starting udp server...")

	for {
		n, addr, err := updLn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		log.Println(n)
		log.Println(addr)

		switch buf[0] {
		case 1:
			//publish
			for _, s := range subscribers.subscribers {
				updLn.WriteTo(buf, s.addr)
			}

		case 2:
			//subscribe
			//var subscriber Subscriber
			log.Printf("subscriber: %s\n", addr)

			//found := false

			sub, ok := subscribers.lookUpSubscriber(addr)

			if ok {
				sub.update()
			} else {
				subscribers.addMember(addr)
			}

			/*
				for _, s := range subscribers.subscribers {
					if s.isSamePort(addr) {
						subscriber = s
						found = true
						break
					}
				}
			*/

			/*
				if !found {
					subscribers.addMember(addr)
				} else {
					log.Printf("update: %s\n", subscriber)
					subscriber.update()
				}
			*/
			log.Println(subscribers)
		}

		/*
			go func() {
				log.Printf("Reciving data: %s from %s", string(buf[:n]), addr.String())

				for i := 0; i < 10; i++ {
					log.Printf("Sending data..")
					message := strconv.Itoa(i)
					updLn.WriteTo([]byte(message), addr)
					log.Printf("Complete Sending data..")
				}
			}()
		*/
	}
}
