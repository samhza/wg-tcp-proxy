package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

func main() {
	var addr netip.AddrPort
	flag.Func("addr", "address to listen on", func(s string) (err error) {
		addr, err = netip.ParseAddrPort(s)
		return
	})
	var target netip.AddrPort
	flag.Func("target", "address to reverse proxy to", func(s string) (err error) {
		target, err = netip.ParseAddrPort(s)
		return
	})
	pubkey := flag.String("pubkey", "", "public key of peer")
	endpoint := flag.String("endpoint", "", "endpoint of peer")
	privkey := flag.String("privkey", "", "our private key")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()
	if !addr.IsValid() || !target.IsValid() || *pubkey == "" || *endpoint == "" || *privkey == "" {
		fmt.Fprintf(os.Stderr, "missing required argument\n")
		flag.Usage()
		os.Exit(1)
	}
	tun, stack, err := netstack.CreateNetTUN([]netip.Addr{addr.Addr()}, nil, 1420)
	if err != nil {
		log.Fatalln("creating tun:", err)
	}
	logger := &device.Logger{Errorf: log.Printf}
	if *verbose {
		logger.Verbosef = log.Printf
	} else {
		logger.Verbosef = func(f string, v ...interface{}) {}
	}
	dev := device.NewDevice(tun, conn.NewStdNetBind(), logger)
	err = dev.IpcSet(fmt.Sprintf(`private_key=%s
public_key=%s
endpoint=%s
allowed_ip=0.0.0.0/0
allowed_ip=::/0
persistent_keepalive_interval=25`,
		convertKey(*privkey), convertKey(*pubkey), *endpoint))
	if err != nil {
		log.Fatalln("setting device:", err)
	}
	err = dev.Up()
	if err != nil {
		log.Fatalln("starting device:", err)
	}
	listener, err := stack.ListenTCP(&net.TCPAddr{IP: net.IP(addr.Addr().AsSlice()), Port: int(addr.Port())})
	taddr := target.String()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accepting:", err)
			continue
		}
		go func() {
			defer conn.Close()
			dst, err := net.Dial("tcp", taddr)
			if err != nil {
				log.Println("dialing:", err)
				return
			}
			defer dst.Close()
			go io.Copy(conn, dst)
			io.Copy(dst, conn)
		}()
	}
}

func convertKey(s string) string {
	p, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatalln("decoding base64:", err)
	}
	return hex.EncodeToString(p)
}
