package main

import (
    "fmt"
    "net"
    "sync"
)

func isPortOpen(host string, port int) bool {
    hostPort := fmt.Sprintf("%s:%d", host, port)

    // resolve address
    tcpAddr, err := net.ResolveTCPAddr("tcp4", hostPort)
    if err != nil {
        return false
    }

    // try to connect
    conn, err := net.Dial("tcp", tcpAddr.String())
    if err != nil {
        return false
    }
    defer conn.Close()

    return true
}

func main() {
    // host addr
    host := "127.0.0.1"
    // ports for scan (portStart-portEnd)
    portStart := 1
    portEnd   := 65535
    // threads count
    threads := 20

    openPorts := []int{}
    mu := sync.Mutex{}
    sem := make(chan bool, threads)

    for port := portStart; port <= portEnd; port++ {
        sem <- true // like wg.Add(1)
        go func(port int) {
            if isPortOpen(host, port) {
                mu.Lock()
                openPorts = append(openPorts, port)
                mu.Unlock()
            }
            <- sem // like wg.Done() & if `sem` have full capacity then wait for goroutines to finish
        }(port)
    }

    // like wg.Wait()
    for i := 0; i < cap(sem); i++ {
        sem <- true
    }

    fmt.Println(openPorts)
}
