package main

import (
	"regexp"

	"github.com/miekg/dns"
)

func main() {
	InitLogging("stderr", debugLevel)

	Infof("Starting DNSProx")

	dnsCache := InitCache(300) /* 5m TTL */
	domains, _ := appConfigs.DNSConfigs["domains"]
	servers, _ := appConfigs.DNSConfigs["servers"]

	dnsProxy := DNSProxy{
		Cache:         &dnsCache,
		domains:       domains.(map[string]interface{}),
		servers:       servers.(map[string]interface{}),
		defaultServer: "8.8.8.8:53",
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:
			m, err := dnsProxy.getResponse(r)
			if err != nil {
				Errorf("Failed lookup for %s with error: %s\n", r, err.Error())
				m.SetReply(r)
				w.WriteMsg(m)
				return
			}
			if len(m.Answer) > 0 {
				pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
				ipAddress := pattern.FindAllString(m.Answer[0].String(), -1)

				if len(ipAddress) > 0 {
					Infof("Lookup for %s with ip %s\n", m.Answer[0].Header().Name, ipAddress[0])
				} else {
					Infof("Lookup for %s with response %s\n", m.Answer[0].Header().Name, m.Answer[0])
				}
			}
			m.SetReply(r)
			w.WriteMsg(m)
		}
	})

	server := &dns.Server{Addr: "127.0.0.1:53", Net: "udp"}
	err := server.ListenAndServe()
	if err != nil {
		Errorf("Failed to start DNSProx: %s\n", err.Error())
	}
}
