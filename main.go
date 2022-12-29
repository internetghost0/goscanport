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
    threads := 20

    openPorts := []int{}
    locking := sync.Mutex{}
    sem := make(chan bool, threads)
    for port := portStart; port <= portEnd; port++ {
        sem <- true
        go func(port int) {
            if isPortOpen(host, port) {
                locking.Lock()
                openPorts = append(openPorts, port)
                locking.Unlock()
            }
            <- sem
        }(port)
    }
    for i := 0; i < cap(sem); i++ {
        sem <- true
    }
    fmt.Println(openPorts)
}
