package main

import (
    "fmt"
    "net"
    "sync"
)

func isPortOpen(host string, port int) bool {
    hostPort := fmt.Sprintf("%s:%d", host, port)
    tcpAddr, err := net.ResolveTCPAddr("tcp4", hostPort)
    if err != nil {
        return false
    }
    conn, err := net.Dial("tcp", tcpAddr.String())
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

func main() {
    host := "127.0.0.1"
    portStart := 1
    portEnd   := 65535

    openPorts := []int{}
    mu := sync.Mutex{}
    var wg sync.WaitGroup
    for port := portStart; port <= portEnd; port++ {
        wg.Add(1)
        go func(port int) {
            if isPortOpen(host, port) {
                mu.Lock()
                openPorts = append(openPorts, port)
                mu.Unlock()
            }
            wg.Done()
        }(port)
    }
    wg.Wait()
    fmt.Println(openPorts)
}
