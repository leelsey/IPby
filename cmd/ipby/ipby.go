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

func checkError(err error) {
	defer func() {
		if recv := recover(); recv != nil {
			fmt.Println("\n"+lstDot, recv)
		}
	}()
	if err != nil {
		panic(clrRed + "Error >> " + err.Error())
	}
}

var (
	appVer        = "0.2"
	appTitle      = clrPurple + "IPby " + clrGrey + "v" + appVer + clrReset
	lstDot        = "  • "
	prvIPTitle    = clrCyan + "Private IP" + clrReset
	pubIPTitle    = clrCyan + "Public IP" + clrReset
	wifiEthernet  = clrBlue + "- WiFi Ethernet" + clrReset
	wiredEthernet = clrBlue + "- Wired Ethernet" + clrReset
	msgIPAddress  = lstDot + clrGreen + "IP Address" + clrReset + ": "
	msgNoSignal   = clrYellow + "No signal" + clrReset
	msgOffline    = msgNoSignal + ": network is turned off."
	msgDisconnect = msgNoSignal + ": internet is disconnected."
	macIF         = "ipconfig"
	macIFWifi     = "en0"
	macIFWired    = "en5"
	macIFAddr     = "getifaddr"
	macIFGetOpt   = "getoption"
	macIFSubMsk   = "subnet_mask"
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
	getTimeout := 10000 * time.Millisecond
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
		fmt.Print(msgIPAddress + string(prvIPWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		fmt.Println(wiredEthernet)
		fmt.Print(msgIPAddress + string(prvIPWireAddr))
	} else {
		fmt.Println(wifiEthernet)
		fmt.Print(msgIPAddress + string(prvIPWiFiAddr))
		fmt.Println(wiredEthernet)
		fmt.Print(msgIPAddress + string(prvIPWireAddr))
	}
}

func getPrvIPMacFull() {
	prvIPWiFi := exec.Command(macIF, macIFAddr, macIFWifi)
	prvIPWire := exec.Command(macIF, macIFAddr, macIFWired)
	netMskWiFi := exec.Command(macIF, macIFGetOpt, macIFWifi, macIFSubMsk)
	netMskWire := exec.Command(macIF, macIFGetOpt, macIFWired, macIFSubMsk)
	routerWiFi := exec.Command(macIF, macIFGetOpt, macIFWifi, macIFRouter)
	routerWire := exec.Command(macIF, macIFGetOpt, macIFWired, macIFRouter)
	prvIPWiFiAddr, _ := prvIPWiFi.Output()
	prvIPWireAddr, _ := prvIPWire.Output()
	netMskWiFiAddr, _ := netMskWiFi.Output()
	netMskWireAddr, _ := netMskWire.Output()
	routerWiFiAddr, _ := routerWiFi.Output()
	routerWireAddr, _ := routerWire.Output()
	if len(prvIPWiFiAddr) == 0 && len(prvIPWireAddr) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else if len(prvIPWireAddr) == 0 {
		fmt.Println(wifiEthernet)
		fmt.Print(msgIPAddress + string(prvIPWiFiAddr) +
			lstDot + "Subnetwork: " + string(netMskWiFiAddr) +
			lstDot + "Net Router: " + string(routerWiFiAddr))
	} else if len(prvIPWiFiAddr) == 0 {
		fmt.Println(wiredEthernet)
		fmt.Print(msgIPAddress + string(prvIPWireAddr) +
			lstDot + "Subnetwork: " + string(netMskWireAddr) +
			lstDot + "Net Router: " + string(routerWireAddr))
	} else {
		fmt.Println(wifiEthernet)
		fmt.Print(msgIPAddress + string(prvIPWiFiAddr) +
			lstDot + "Subnetwork: " + string(netMskWiFiAddr) +
			lstDot + "Net Router: " + string(routerWiFiAddr))
		fmt.Println(wiredEthernet)
		fmt.Print(msgIPAddress + string(prvIPWireAddr) +
			lstDot + "Subnetwork: " + string(netMskWireAddr) +
			lstDot + "Net Router: " + string(routerWireAddr))
	}
}

func getPrvIPLinux() {
	prvIP := exec.Command(linuxIF, linuxIFOpt)
	prvIPAddr, _ := prvIP.Output()
	if len(prvIPAddr) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else {
		fmt.Print(msgIPAddress + string(prvIPAddr))
	}
}

func getPrvIPWindows() {
	var psList []*exec.Cmd
	psList = append(psList, exec.Command("powershell", "/C", winIF))
	psList = append(psList, exec.Command("powershell", "/C", "$Input | ", winIFOpt, winIFIPv4))
	var prvIP []byte
	for i, s := range psList {
		if i > 0 {
			input, err := s.StdinPipe()
			checkError(err)
			//if err != nil {
			//	panic(err)
			//}
			go func(write io.WriteCloser, data []byte) {
				_, errWrite := write.Write(data)
				checkError(errWrite)
				//if err != nil {
				//	return
				//}
				ereClose := write.Close()
				checkError(ereClose)
				//if err != nil {
				//	return
				//}
			}(input, prvIP)
		}
		var err error
		prvIP, err = s.CombinedOutput()
		checkError(err)
		//if err != nil {
		//	panic(err)
		//}
	}
	prvIPList := strings.Split(string(prvIP), ":")
	if len(prvIPList) == 0 {
		fmt.Println(lstDot + msgOffline)
	} else {
		for listNum, prvIPAddr := range prvIPList {
			if listNum >= 1 {
				fmt.Println(msgIPAddress + prvIPAddr[1:15])
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
		var pubIPv4 = getPubIP()
		var pubIPv6 = getPubIP64()

		if pubIPv4 == pubIPv6 {
			fmt.Println(msgIPAddress + pubIPv4)
		} else {
			fmt.Println(lstDot + "IP Address v4: " + pubIPv4)
			fmt.Println(lstDot + "IP Address v6: " + pubIPv6)
		}
	} else {
		fmt.Println(lstDot + msgDisconnect)
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
			getPrvIPWindows()
		}
		fmt.Println(pubIPTitle)
		if checkNetStatus() == true {
			fmt.Println(msgIPAddress + getPubIP64())
		} else {
			fmt.Println(lstDot + msgDisconnect)
		}
	} else {
		switch os.Args[1] {
		case "version":
			err := verOpt.Parse(os.Args[1:])
			checkError(err)
			fmt.Println(appTitle)
		case "all":
			err := allOpt.Parse(os.Args[1:])
			checkError(err)
			fmt.Println(prvIPTitle)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			case "windows":
				getPrvIPWindows()
			}
			fmt.Println(pubIPTitle)
			pubIPAll()
		case "public":
			err := pubOpt.Parse(os.Args[1:])
			checkError(err)
			fmt.Println(pubIPTitle)
			pubIPAll()
		case "private":
			err := prvOpt.Parse(os.Args[1:])
			checkError(err)
			fmt.Println(prvIPTitle)
			switch runtime.GOOS {
			case "darwin":
				getPrvIPMacFull()
			case "linux":
				getPrvIPLinux()
			case "windows":
				getPrvIPWindows()
			}
		case "help":
			err := prvOpt.Parse(os.Args[1:])
			checkError(err)
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
