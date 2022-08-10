package main

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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

func mapNoneEmpty[R netip.Addr | string](data []string, f func(string) R) []R {

	mapped := make([]R, len(data))
	i := 0
	for _, e := range data {
		if e != "" {
			mapped[i] = f(e)
			i++
		}
	}

	return mapped[:i]
}

func DecodeKey(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Fatalf("Error decoding either preshared, private or public key from base64: '%v'", err)
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

func tern(prefix string, s string) string {
	if s == "" {
		return ""
	}
	return prefix + "=" + s
}

func GenerateConfig(cfg *Cfg) string {
	// https://www.wireguard.com/xplatform/#configuration-protocol
	// todo add all possible settings
	//todo set array cap
	commands := []string{
		tern("private_key", DecodeKey(cfg.PrivateKey)),     //todo required
		tern("public_key", DecodeKey(cfg.PublicKey)),       //todo required
		tern("preshared_key", DecodeKey(cfg.PresharedKey)), //todo optional
		tern("endpoint", cfg.Endpoint),                     //todo required
	}
	if cfg.PersistentKeepalive != 0 {
		commands = append(commands, "persistent_keepalive_interval="+strconv.FormatUint(cfg.PersistentKeepalive, 10))
	}
	if cfg.ListenPort != 0 {
		commands = append(commands, "listen_port="+strconv.FormatUint(cfg.ListenPort, 10))
	}
	commands = append(commands, mapNoneEmpty(cfg.AllowedIPs, func(addr string) string {
		return "allowed_ip=" + addr
	})...)
	return strings.Join(commands, "\n")
}

func main() {
	lg := log.New(os.Stdout, "", log.LstdFlags)

	cli, err := ParseFlags()
	if err != nil {
		lg.Fatalf("Fail to parse cli args: %v", err)
	}

	cfg := &Cfg{}
	err = ini.MapTo(cfg, cli.ConfigFile) //todo allow load from env
	if err != nil {
		lg.Fatalf("Fail to load file: %v", err)
	}

	lg.Println("Creating TUN")

	if cfg.MTU == 0 {
		cfg.MTU = 1420
	}
	tun, tnet, err := netstack.CreateNetTUN(mapNoneEmpty(cfg.Address, mustParseCIDR), mapNoneEmpty(cfg.DNS, netip.MustParseAddr), cfg.MTU)
	if err != nil {
		lg.Fatalf("Error creating TUN '%v'", err)
	}

	lg.Println("Setting up wireguard connection")

	dev := device.NewDevice(tun, conn.NewStdNetBind(), device.NewLogger(device.LogLevelError, "")) //todo logger
	err = dev.IpcSet(GenerateConfig(cfg))
	if err != nil {
		lg.Fatal(err)
	}
	err = dev.Up()
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println("Starting proxy server on port", cli.Port)

	proxyUrl, _ := url.Parse("http://127.0.0.1:" + cli.SPort())
	proxy := &goproxy.ProxyHttpServer{
		Verbose: true,
		Logger:  lg,
		NonproxyHandler: http.HandlerFunc(ProxyHandler(http.Client{
			Timeout:   4 * time.Minute,
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		})),
		Tr: &http.Transport{
			Dial:        tnet.Dial,
			DialContext: tnet.DialContext,
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		ConnectDial: tnet.Dial,
	}

	//debug.ReadBuildInfo();

	/*if *silent {
		proxy.Logger.SetOutput(io.Discard)
	} else {
		proxy.Verbose = *verbose
	}*/

	go func() {
		lg.Fatal(http.ListenAndServe(":"+cli.SPort(), proxy))
	}()
	res, err := http.Head("http://127.0.0.1:" + cli.SPort() + "/health") //todo check some host like google
	if err != nil || res == nil {
		lg.Fatalf("Can't check if proxy server started '%v'", err)
	} else if res.StatusCode == 204 {
		lg.Println("Server started")
	} else {
		lg.Fatalf("Can't check if proxy server started '%v'", res.Status)
	}
	if !cli.StartAndExit {
		select {}
	}
}
