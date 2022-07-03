package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	appver           = "0.1"
	lstdot           = "  • "
	titlefnt         = color.New(color.FgGreen, color.Bold)
	prvipt           = "Private IP"
	pubipt           = "Public IP"
	offline          = "Network is turned off"
	disconnet        = "Internet disconnected"
	ifmac            = "ipconfig" // Private IP for macOS
	ifmac_netwifi    = "en0"
	ifmac_netwire    = "en5"
	ifmac_ifaddr     = "getifaddr"
	ifmac_getoption  = "getoption"
	ifmac_subnetmask = "subnet_mask"
	ifmac_router     = "router"
	iflinux          = "hostname" // Private IP for Linux
	iflinux_opt      = "-I"
)

func getPrvIPMacSimple() {
	prvIPWiFi := exec.Command(ifmac, ifmac_ifaddr, ifmac_netwifi)
	prvIPWire := exec.Command(ifmac, ifmac_ifaddr, ifmac_netwire)
	prvIPWiFiAddr, _ := prvIPWiFi.Output()
	prvIPWireAddr, _ := prvIPWire.Output()
	if len(prvIPWiFiAddr) == 0 && len(prvIPWireAddr) == 0 {
		fmt.Println(lstdot + offline)
	} else if len(prvIPWireAddr) == 0 {
		color.Cyan("- WiFi Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		color.Cyan("- Wired Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWireAddr))
	} else {
		color.Cyan("- WiFi Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWiFiAddr))
		color.Cyan("- Wired Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWireAddr))
	}
}

func getPrvIPMacFull() {
	prvIPWiFi := exec.Command(ifmac, ifmac_ifaddr, ifmac_netwifi)
	prvIPWire := exec.Command(ifmac, ifmac_ifaddr, ifmac_netwire)
	netMskWiFi := exec.Command(ifmac, ifmac_getoption, ifmac_netwifi, ifmac_subnetmask)
	netMskWire := exec.Command(ifmac, ifmac_getoption, ifmac_netwire, ifmac_subnetmask)
	routerWiFi := exec.Command(ifmac, ifmac_getoption, ifmac_netwifi, ifmac_router)
	routerWire := exec.Command(ifmac, ifmac_getoption, ifmac_netwire, ifmac_router)
	prvIPWiFiAddr, _ := prvIPWiFi.Output()
	prvIPWireAddr, _ := prvIPWire.Output()
	netmskWiFiAddr, _ := netMskWiFi.Output()
	netmskWireAddr, _ := netMskWire.Output()
	routerWiFiAddr, _ := routerWiFi.Output()
	routerWireAddr, _ := routerWire.Output()
	if len(prvIPWiFiAddr) == 0 && len(prvIPWireAddr) == 0 {
		fmt.Println(lstdot + offline)
	} else if len(prvIPWireAddr) == 0 {
		color.Cyan("- WiFi Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWiFiAddr) +
			lstdot + "Subnetwork: " + string(netmskWiFiAddr) +
			lstdot + "Net Router: " + string(routerWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		color.Cyan("- Wired Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWireAddr) +
			lstdot + "Subnetwork: " + string(netmskWireAddr) +
			lstdot + "Net Router: " + string(routerWireAddr))
	} else {
		color.Cyan("- WiFi Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWiFiAddr) +
			lstdot + "Subnetwork: " + string(netmskWiFiAddr) +
			lstdot + "Net Router: " + string(routerWiFiAddr))
		color.Cyan("- Wired Ethernet")
		fmt.Print(lstdot + "IP Address: " + string(prvIPWireAddr) +
			lstdot + "Subnetwork: " + string(netmskWireAddr) +
			lstdot + "Net Router: " + string(routerWireAddr))
	}
}

func getPrvIPLinux() {
	prvip := exec.Command(iflinux, iflinux_opt)
	prvipAddr, _ := prvip.Output()
	if len(prvipAddr) == 0 {
		fmt.Println(lstdot + offline)
	} else {
		fmt.Print(lstdot + "IP Address: " + string(prvipAddr))
	}
}

//func getPrvIPWIndows() {
//	return
//}

func getPubIP() string {
	ipv4, _ := http.Get("https://api.ipify.org")
	ipv4addr, _ := ioutil.ReadAll(ipv4.Body)
	return string(ipv4addr)
}

func getPubIP64() string {
	ipv6, _ := http.Get("https://api64.ipify.org")
	ipv6addr, _ := ioutil.ReadAll(ipv6.Body)
	return string(ipv6addr)
}

func pubIPAll() {
	if checkNetStatus() == true {
		var pubipv4 string = getPubIP()
		var pubipv6 string = getPubIP64()
		if pubipv4 == pubipv6 {
			fmt.Println(lstdot + "IP Address: " + pubipv4)
		} else {
			fmt.Println(lstdot + "IP Address v4: " + pubipv4)
			fmt.Println(lstdot + "IP Address v6: " + pubipv6)
		}
	} else {
		fmt.Println(lstdot + disconnet)
	}
}

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

func main() {
	verOpt := flag.NewFlagSet("version", flag.ExitOnError)
	allOpt := flag.NewFlagSet("all", flag.ExitOnError)
	pubOpt := flag.NewFlagSet("public", flag.ExitOnError)
	prvOpt := flag.NewFlagSet("private", flag.ExitOnError)
	if len(os.Args) == 1 {
		titlefnt.Println(prvipt)
		switch runtime.GOOS {
		case "darwin":
			getPrvIPMacSimple()
		case "linux":
			getPrvIPLinux()
		}
		titlefnt.Println(pubipt)
		if checkNetStatus() == true {
			fmt.Println(lstdot + "IP Address: " + getPubIP64())
		} else {
			fmt.Println(lstdot + disconnet)
		}
	} else {
		switch os.Args[1] {
		case "version":
			verOpt.Parse(os.Args[1:])
			fmt.Println(lstdot + "Version: " + appver)
		case "all":
			allOpt.Parse(os.Args[1:])
			titlefnt.Println(prvipt)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			}
			titlefnt.Println(pubipt)
			pubIPAll()
		case "public":
			pubOpt.Parse(os.Args[1:])
			titlefnt.Println(pubipt)
			pubIPAll()
		case "private":
			prvOpt.Parse(os.Args[1:])
			titlefnt.Println(prvipt)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			}
		default:
			fmt.Println(lstdot + "Usage: ipby <command>\n" +
				lstdot + "It can use in <command> that one of version, all, public, private")
		}
	}
}
