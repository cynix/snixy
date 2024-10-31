package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/inetaf/tcpproxy"
)

func main() {
	version := flag.Bool("version", false, "show version and exit")
	cf := flag.String("config", "/usr/local/etc/snixy.yaml", "config  file")
	flag.Parse()

	if (*version) {
		fmt.Println(Version)
		return
	}

	var config Config

	if err := config.Load(*cf); err != nil {
		log.Fatal(err)
	}

	var p tcpproxy.Proxy
	defer p.Close()

	for _, route := range config.Routes {
		if len(route.SNI) > 0 {
			for _, sni := range route.SNI {
				p.AddSNIRoute(route.Listen, sni, tcpproxy.To(route.Dial))
			}
		} else {
			p.AddRoute(route.Listen, tcpproxy.To(route.Dial))
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := p.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sig
}

var (
	Version string = "dev"
)
