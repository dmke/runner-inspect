package main

import (
	"log"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const nonExistingNameserver = ":7053"

func main() {
	log.SetFlags(log.Lmicroseconds)

	m := new(dns.Msg)
	m.SetQuestion("mail.google.com.", dns.TypeSOA)
	m.SetEdns0(4096, false)

	wg := sync.WaitGroup{}
	for _, proto := range []string{
		"udp", "udp4", "udp6",
		"tcp", "tcp4", "tcp6",
	} {
		wg.Add(1)
		go func(proto string) {
			try(proto, m)
			wg.Done()
		}(proto)
	}

	wg.Wait()
	log.Println("FINISH")
}

func try(proto string, m *dns.Msg) {
	client := dns.Client{Net: proto, Timeout: time.Second}

	log.Printf("START proto: %s", proto)
	response, ttl, err := client.Exchange(m, nonExistingNameserver)
	if err == nil {
		log.Printf("OK    proto: %s, ttl: %v, result: %s", proto, ttl, response.String())
	} else {
		log.Printf("ERROR proto: %s, error: %v", proto, err)
	}
}
