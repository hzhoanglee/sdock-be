package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ScanDevice struct {
	IP           string `json:"ip"`
	ScanDeviceID string `json:"device_id"`
}

func ScanForDevices() []ScanDevice {
	devices := make([]ScanDevice, 0)
	subnetIP, err := getSubnetInfo()
	if err != nil {
		fmt.Println(err)
		return devices
	}

	networkPrefix := subnetIP[:strings.LastIndex(subnetIP, ".")]
	var wg sync.WaitGroup
	results := make(chan string)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", networkPrefix, i)
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", ip), 5000*time.Millisecond)
			if err == nil {
				conn.Close()
				results <- ip
			}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for server := range results {
		fmt.Println("Found device at", server)
		resp, err := http.Get(fmt.Sprintf("http://%s/available", server))
		if err == nil && resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				devices = append(devices, ScanDevice{
					IP:           server,
					ScanDeviceID: string(body),
				})
			}
			resp.Body.Close()
		}
	}

	return devices
}

func getSubnetInfo() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("could not find a valid IPv4 subnet")
}
