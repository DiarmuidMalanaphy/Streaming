package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func listen(reqChan chan<- networkData) {

	addr, err := net.ResolveUDPAddr("udp", ":8000") // Listen on all interfaces
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	buffer := make([]byte, 1424) //Arbitrary

	publicIp, err := getPublicIP()
	if err != nil {
		fmt.Println("Error getting public IP:", err)

	}
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP address:", err)
		return
	}

	fmt.Printf("UDP server listening on port 8000 and Global IP - %s\n", publicIp)
	fmt.Printf("Server local IP is - %s\n", localIP)
	for {

		n, remoteAddr, err := conn.ReadFromUDP(buffer)

		req, err := deserialiseRequest(buffer[:n])

		if err != nil {
			// fmt.Println(""err)
			continue
		}

		// We then relay this information back to main through the requests channel
		// Useful way for goroutines to communicate.
		reqChan <- networkData{Request: req, Addr: remoteAddr}
	}
}

func sendUDP(address string, data []byte) error {
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		return err
	}

	// Establish a UDP connection
	conn, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the data

	_, err = conn.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // Interface down or loopback
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Skip loopback and undefined addresses
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			ipStr := ip.String()

			// Check if the IP address is from a common virtual network range
			if !strings.HasPrefix(ipStr, "192.168.56.") {
				return ipStr, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find local IP address")
}
