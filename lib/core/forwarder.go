package core

import (
    "io"
    "net"
    "time"
    "strconv"
    "log"
)

type Info struct {
    SourceAddr  string
    SourcePort  int
    TargetAddr  string
    TargetPort  int
    Protocol    string
}

func copyIO(source, target net.Conn) {
    if source == nil || target == nil {
        if source != nil {
            source.Close()
        }
        if target != nil {
            target.Close()
        }
        return
    }
    
    defer source.Close()
    defer target.Close()
    
    sourceAddr := source.RemoteAddr().String()
    targetAddr := target.RemoteAddr().String()
    
    written, err := io.Copy(target, source)
    if err != nil {
        if err != io.EOF {
            log.Printf("Error copying data from %s to %s: %v", sourceAddr, targetAddr, err)
        }
    }
    
    log.Printf("Connection closed between %s and %s, bytes transferred: %d", 
        sourceAddr, targetAddr, written)
}

func Forward(info Info) (net.Listener, error) {
    // If source address is empty, use "0.0.0.0" to listen on all interfaces
    sourceAddr := info.SourceAddr
    if sourceAddr == "" {
        sourceAddr = "0.0.0.0"
    }
    
    listenAddr := net.JoinHostPort(sourceAddr, strconv.Itoa(info.SourcePort))
    log.Printf("Starting forward from %s to %s:%d using %s", 
        listenAddr, info.TargetAddr, info.TargetPort, info.Protocol)
    
    ln, err := net.Listen(info.Protocol, listenAddr)
    if err != nil {
        log.Printf("Failed to listen on %s: %v", listenAddr, err)
        return nil, err
    }
    
    log.Printf("Successfully listening on %s", listenAddr)
    
    go func() {
        for {
            source, err := ln.Accept()
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Err == net.ErrClosed {
                    log.Printf("Listener closed, stopping forward from %s", listenAddr)
                    return
                }
                log.Printf("Accept error on %s: %v", listenAddr, err)
                continue
            }
            
            clientAddr := source.RemoteAddr().String()
            log.Printf("New connection from %s to %s", clientAddr, listenAddr)
            
            go func() {
                targetAddr := net.JoinHostPort(info.TargetAddr, strconv.Itoa(info.TargetPort))
                log.Printf("Forwarding connection from %s to %s", clientAddr, targetAddr)
                
                target, err := net.Dial(info.Protocol, targetAddr)
                if err != nil {
                    log.Printf("Failed to connect to target %s: %v", targetAddr, err)
                    source.Close()
                    return
                }
                
                if target == nil {
                    log.Printf("Failed to establish target connection to %s", targetAddr)
                    source.Close()
                    return
                }
                
                log.Printf("Successfully established connection to target %s", targetAddr)
                
                go copyIO(source, target)
                go copyIO(target, source)
            }()
        }
    }()
    return ln, err
}

func Telnet(protocol, ip string, port, timeout int) (bool, error) {
    address := net.JoinHostPort(ip, strconv.Itoa(port))
    log.Printf("Testing connection to %s://%s with timeout %dms", protocol, address, timeout)
    
    conn, err := net.DialTimeout(protocol, address, time.Duration(timeout) * time.Millisecond)
    if err != nil {
        log.Printf("Connection test failed to %s://%s: %v", protocol, address, err)
        return false, err
    }
    
    log.Printf("Connection test successful to %s://%s", protocol, address)
    conn.Close()
    return true, nil
}
