package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

type dnsHandler struct{}

// based on https://reintech.io/blog/implementing-a-dns-server-in-go

func resolve(domain string, qtype uint16) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	rr, err := dns.NewRR(fmt.Sprintf("%s PTR this.is.a.test", domain))

	if domain == "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.4.b.0.8.2.0.0.0.7.4.0.1.0.0.2.ip6.arpa." {
		rr, err = dns.NewRR(fmt.Sprintf("%s PTR one.two.three.four", domain))
	}

	if domain == "2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.4.b.0.8.2.0.0.0.7.4.0.1.0.0.2.ip6.arpa." {
		rr, err = dns.NewRR(fmt.Sprintf("%s PTR two.three.four", domain))
	}

	if domain == "3.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.4.b.0.8.2.0.0.0.7.4.0.1.0.0.2.ip6.arpa." {
		rr, err = dns.NewRR(fmt.Sprintf("%s PTR three.four", domain))
	}

	if err != nil {
		log.Fatal(err)
	}
	rrArr := []dns.RR{rr}

	return rrArr, nil
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		fmt.Printf("Received query: %s\n", question.Name)
		answers, err := resolve(question.Name, question.Qtype)
		if err != nil {
			log.Fatal(err)
		}
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
}

func main() {
	handler := new(dnsHandler)
	server := &dns.Server{
		Addr:      ":5353",
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port 53")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
