package core

import (
    "io"
    "net"
    "time"
    "strconv"
)

type Info struct {
    SourceAddr  string
    SourcePort  int
    TargetAddr  string
    TargetPort  int
    Protocol    string
}

func copyIO(source, target net.Conn) {
    if source != nil {
        defer source.Close()
    }
    if target != nil {
        defer target.Close()
    }
    if source != nil && target != nil {
        io.Copy(source, target)
    }
}

func Forward(info Info) (net.Listener, error) {
    ln, err := net.Listen(info.Protocol, net.JoinHostPort(info.SourceAddr, strconv.Itoa(info.SourcePort)))
    if err == nil {
        go func() {
            for {
                source, err := ln.Accept()
                if err != nil {
                    if opErr, ok := err.(*net.OpError); ok && opErr.Err == net.ErrClosed {
                        return
                    }
                }
                go func() {
                    target, _ := net.Dial(info.Protocol, net.JoinHostPort(info.TargetAddr, strconv.Itoa(info.TargetPort)))
                    go copyIO(source, target)
                    go copyIO(target, source)
                }()
            }
        }()
    }
    return ln, err
}

func Telnet(protocol, ip string, port, timeout int) (bool, error) {
    conn, err := net.DialTimeout(protocol, net.JoinHostPort(ip, strconv.Itoa(port)), time.Duration(timeout) * time.Millisecond)
    if err != nil {
        return false, err
    }
    conn.Close()
    return true, nil
}
