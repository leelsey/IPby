package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	appVer        = "0.2"
	appTitle      = clrPurple + "IPBy " + clrGrey + "v" + appVer + clrReset
	lstDot        = "  • "
	prvIPTitle    = clrCyan + "Private IP" + clrReset
	pubIPTitle    = clrCyan + "Public IP" + clrReset
	wifiEthernet  = clrBlue + "- WiFi Ethernet" + clrReset
	wiredEthernet = clrBlue + "- Wired Ethernet" + clrReset
	msgNoSignal   = clrYellow + "No signal" + clrReset
	msgOffline    = msgNoSignal + ": network is turned off."
	msgDisconnet  = msgNoSignal + ": internet is disconnected."
	macIF         = "ipconfig"
	macIFWifi     = "en0"
	macIFWired    = "en5"
	macIFAddr     = "getifaddr"
	macIFGetOpt   = "getoption"
	macIFSBW      = "subnet_mask"
	macIFRouter   = "router"
	linuxIF       = "hostname"
	linuxIFOpt    = "-I"
	winIF         = "ipconfig"
	winIFOpt      = "findstr"
	winIFIPv4     = "IPv4"
	//winIFIPv6     = "IPv6"
	clrReset  = "\033[0m"
	clrRed    = "\033[31m"
	clrGreen  = "\033[32m"
	clrYellow = "\033[33m"
	clrBlue   = "\033[34m"
	clrPurple = "\033[35m"
	clrCyan   = "\033[36m"
	clrGrey   = "\033[37m"
	//clrWhite  = "\033[97m"
	//clrBlack  = "\033[30m"
)

func checkNetStatus() bool {
	getTimeout := time.Duration(10000 * time.Millisecond)
	client := http.Client{
		Timeout: getTimeout,
	}
	_, err := client.Get("https://9.9.9.9")
	if err != nil {
		return false
	}
	return true
}

func getPrvIPMacSimple() {
	prvIPWiFi := exec.Command(macIF, macIFAddr, macIFWifi)
	prvIPWire := exec.Command(macIF, macIFAddr, macIFWired)
	prvIPWiFiAddr, _ := prvIPWiFi.Output()
	prvIPWireAddr, _ := prvIPWire.Output()
	if len(prvIPWiFiAddr) == 0 && len(prvIPWireAddr) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else if len(prvIPWireAddr) == 0 {
		fmt.Println(wifiEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		fmt.Println(wiredEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWireAddr))
	} else {
		fmt.Println(wifiEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWiFiAddr))
		fmt.Println(wiredEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWireAddr))
	}
}

func getPrvIPMacFull() {
	prvIPWiFi := exec.Command(macIF, macIFAddr, macIFWifi)
	prvIPWire := exec.Command(macIF, macIFAddr, macIFWired)
	netMskWiFi := exec.Command(macIF, macIFGetOpt, macIFWifi, macIFSBW)
	netMskWire := exec.Command(macIF, macIFGetOpt, macIFWired, macIFSBW)
	routerWiFi := exec.Command(macIF, macIFGetOpt, macIFWifi, macIFRouter)
	routerWire := exec.Command(macIF, macIFGetOpt, macIFWired, macIFRouter)
	prvIPWiFiAddr, _ := prvIPWiFi.Output()
	prvIPWireAddr, _ := prvIPWire.Output()
	netmskWiFiAddr, _ := netMskWiFi.Output()
	netmskWireAddr, _ := netMskWire.Output()
	routerWiFiAddr, _ := routerWiFi.Output()
	routerWireAddr, _ := routerWire.Output()
	if len(prvIPWiFiAddr) == 0 && len(prvIPWireAddr) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else if len(prvIPWireAddr) == 0 {
		fmt.Println(wifiEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWiFiAddr) +
			lstDot + "Subnetwork: " + string(netmskWiFiAddr) +
			lstDot + "Net Router: " + string(routerWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		fmt.Println(wiredEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWireAddr) +
			lstDot + "Subnetwork: " + string(netmskWireAddr) +
			lstDot + "Net Router: " + string(routerWireAddr))
	} else {
		fmt.Println(wifiEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWiFiAddr) +
			lstDot + "Subnetwork: " + string(netmskWiFiAddr) +
			lstDot + "Net Router: " + string(routerWiFiAddr))
		fmt.Println(wiredEthernet)
		fmt.Print(lstDot + "IP Address: " + string(prvIPWireAddr) +
			lstDot + "Subnetwork: " + string(netmskWireAddr) +
			lstDot + "Net Router: " + string(routerWireAddr))
	}
}

func getPrvIPLinux() {
	prvip := exec.Command(linuxIF, linuxIFOpt)
	prvipAddr, _ := prvip.Output()
	if len(prvipAddr) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else {
		fmt.Print(lstDot + "IP Address: " + string(prvipAddr))
	}
}

func getPrvIPWIndows() {
	var psList []*exec.Cmd
	psList = append(psList, exec.Command("powershell", "/C", winIF))
	psList = append(psList, exec.Command("powershell", "/C", "$Input | ", winIFOpt, winIFIPv4))
	var prvip []byte
	for i, s := range psList {
		if i > 0 {
			input, err := s.StdinPipe()
			if err != nil {
				panic(err)
			}
			go func(write io.WriteCloser, data []byte) {
				write.Write(data)
				write.Close()
			}(input, prvip)
		}
		var err error
		prvip, err = s.CombinedOutput()
		if err != nil {
			panic(err)
		}
	}
	prvipList := strings.Split(string(prvip), ":")
	if len(prvipList) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else {
		for listnum, prvipAddr := range prvipList {
			if listnum >= 1 {
				fmt.Println(lstDot + "IP Address: " + prvipAddr[1:15])
			}
		}
	}
}

func getPubIP() string {
	ipv4, _ := http.Get("https://api.ipify.org")
	ipv4addr, _ := io.ReadAll(ipv4.Body)
	return string(ipv4addr)
}

func getPubIP64() string {
	ipv6, _ := http.Get("https://api64.ipify.org")
	ipv6addr, _ := io.ReadAll(ipv6.Body)
	return string(ipv6addr)
}

func pubIPAll() {
	if checkNetStatus() == true {
		var pubipv4 string = getPubIP()
		var pubipv6 string = getPubIP64()
		if pubipv4 == pubipv6 {
			fmt.Println(lstDot + "IP Address: " + pubipv4)
		} else {
			fmt.Println(lstDot + "IP Address v4: " + pubipv4)
			fmt.Println(lstDot + "IP Address v6: " + pubipv6)
		}
	} else {
		fmt.Println(lstDot + msgDisconnet)
	}
}

func main() {
	verOpt := flag.NewFlagSet("version", flag.ExitOnError)
	allOpt := flag.NewFlagSet("all", flag.ExitOnError)
	pubOpt := flag.NewFlagSet("public", flag.ExitOnError)
	prvOpt := flag.NewFlagSet("private", flag.ExitOnError)
	if len(os.Args) == 1 {
		fmt.Println(prvIPTitle)
		switch runtime.GOOS {
		case "darwin":
			getPrvIPMacSimple()
		case "linux":
			getPrvIPLinux()
		case "windows":
			getPrvIPWIndows()
		}
		fmt.Println(pubIPTitle)
		if checkNetStatus() == true {
			fmt.Println(lstDot + "IP Address: " + getPubIP64())
		} else {
			fmt.Println(lstDot + msgDisconnet)
		}
	} else {
		switch os.Args[1] {
		case "version":
			verOpt.Parse(os.Args[1:])
			fmt.Println(appTitle)
		case "all":
			allOpt.Parse(os.Args[1:])
			fmt.Println(prvIPTitle)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			case "windows":
				getPrvIPWIndows()
			}
			fmt.Println(pubIPTitle)
			pubIPAll()
		case "public":
			pubOpt.Parse(os.Args[1:])
			fmt.Println(pubIPTitle)
			pubIPAll()
		case "private":
			prvOpt.Parse(os.Args[1:])
			fmt.Println(prvIPTitle)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			case "windows":
				getPrvIPWIndows()
			}
		case "help":
			prvOpt.Parse(os.Args[1:])
			fmt.Println(appTitle)
			fmt.Println(lstDot + "Usage: ipby <command>\n" +
				lstDot + "It can use in <command> that one of version, all, public, private")
		default:
			fmt.Println(clrRed + "Wrong usage\n" + clrReset)
			fmt.Println(lstDot + "Usage: ipby <command>\n" +
				lstDot + "It can use in <command> that one of version, all, public, private")
		}
	}
}
