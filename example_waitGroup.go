package main

import (
    "fmt"
    "net"
    "sync"
)

func main() {
    // host to scan
    host := "127.0.0.1"
    // ports to scan
    portStart := 1
    portEnd   := 65535

    // result
    openPorts := []int{}
    // for locking/unlocking openPorts variable
    mu := sync.Mutex{}
    // for managing goroutines
    var wg sync.WaitGroup

    for port := portStart; port <= portEnd; port++ {
        // increase task counter (task is going to start)
        wg.Add(1)

        // create a new thread for scan a port
        go func(port int) {
            // if port is open
            if isPortOpen(host, port) {
                mu.Lock()
                // append port to result
                openPorts = append(openPorts, port)
                mu.Unlock()
            }
            //decrease task counter (task is done)
            wg.Done()
        }(port)
    }

    // wait for tasks to done
    wg.Wait()
    // print result
    fmt.Println(openPorts)
}

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
