package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"log"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"gopkg.in/elazarl/goproxy.v1"
	"gopkg.in/ini.v1"
)

func mustParseCIDR(s string) netip.Addr {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	nIp, err := netip.ParseAddr(ip.String())
	if err != nil {
		panic(err)
	}
	return nIp
}

func withoutEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func mapA[R netip.Addr | string](data []string, f func(string) R) []R {

	mapped := make([]R, len(data))

	for i, e := range data {
		mapped[i] = f(e)
	}

	return mapped
}

func DecodeKey(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Panicf("Error decoding either preshared, private or public key from base64: '%v'", err)
	}
	return hex.EncodeToString(b)
}

type Peer struct {
	PublicKey           string
	PresharedKey        string
	AllowedIPs          []string
	Endpoint            string
	PersistentKeepalive uint64
}

type Interface struct {
	PrivateKey string
	ListenPort uint64
	MTU        int
	Address    []string
	DNS        []string
}

type Cfg struct {
	Interface
	Peer
}

func formatAllowed(addr string) string {
	return "allowed_ip=" + addr
}

func tern(prefix string, s string) string {
	if s == "" {
		return ""
	}
	return prefix + "=" + s
}

func GenerateConfig(cfg *Cfg) string {
	commands := []string{
		tern("private_key", DecodeKey(cfg.PrivateKey)),
		tern("public_key", DecodeKey(cfg.PublicKey)),
		tern("preshared_key", DecodeKey(cfg.PresharedKey)),
		tern("endpoint", cfg.Endpoint),
	}
	if cfg.PersistentKeepalive != 0 {
		commands = append(commands, tern("persistent_keepalive_interval", strconv.FormatUint(cfg.PersistentKeepalive, 10)))
	}
	if cfg.ListenPort != 0 {
		commands = append(commands, tern("listen_port", strconv.FormatUint(cfg.ListenPort, 10)))
	}
	commands = append(commands, mapA(withoutEmpty(cfg.AllowedIPs), formatAllowed)...)
	return strings.Join(commands, "\n")
}

func main() {

	portPointer := flag.Uint64("p", 8087, "Proxy port")
	silent := flag.Bool("s", false, "Silent mode")
	verbose := flag.Bool("v", false, "Log information on each request sent to the proxy")
	flag.Parse()

	cfg := new(Cfg)
	err := ini.MapTo(cfg, flag.Args()[0])
	if err != nil {
		log.Panicf("Fail to load file: %v", err)
	}

	//ip, _, err := net.ParseCIDR(cfg.Address[0])
	//fmt.Printf("%+v\n", ip)
	if !*silent {
		log.Printf("Creating TUN")
	}
	if cfg.MTU == 0 {
		cfg.MTU = 1420
	}
	tun, tnet, err := netstack.CreateNetTUN(mapA(withoutEmpty(cfg.Address), mustParseCIDR), mapA(withoutEmpty(cfg.DNS), netip.MustParseAddr), cfg.MTU)
	if err != nil {
		log.Panicf("Error creating TUN '%v'", err)
	}

	if !*silent {
		log.Printf("Setting up wireguard connection")
	}
	dev := device.NewDevice(tun, conn.NewStdNetBind(), device.NewLogger(device.LogLevelError, ""))
	dev.IpcSet(GenerateConfig(cfg))
	dev.Up()

	port := strconv.FormatUint(*portPointer, 10)
	if !*silent {
		log.Printf("Starting proxy server on port %v", port)
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.ConnectDial = tnet.Dial
	proxy.Tr.Dial = tnet.Dial
	if !*silent {
		proxy.Verbose = *verbose
	}
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}
