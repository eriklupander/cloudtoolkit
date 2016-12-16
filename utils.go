package cloudtoolkit

import (
        "net"
)

func DumpDNS() {
        arr := [...]string{"configserver", "zipkin", "accountservice", "imageservice", "rabbitmq"}
        for _, val := range arr {
                Log.Println("Looking up hostname: " + val)
                ips, err := net.LookupHost(val)
                if err != nil {
                        Log.Println("Lookup for '" + val + "' returned error: " + err.Error())
                } else {
                        for _, ip := range ips {
                                Log.Println(val + " => " + ip)
                        }
                }
        }
}

func GetLocalIP() string {
        addrs, err := net.InterfaceAddrs()
        if err != nil {
                return ""
        }
        for _, address := range addrs {
                // check the address type and if it is not a loopback the display it
                if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                        if ipnet.IP.To4() != nil {
                                return ipnet.IP.String()
                        }
                }
        }
        panic("Unable to determine local IP address (non loopback). Exiting.")
}

