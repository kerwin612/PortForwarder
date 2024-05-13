package core

import (
    "io"
    "net"
    "time"
    "testing"
)

func Test_Forward(t *testing.T) {
    source, target := net.Pipe()

    info := Info{
        SourceAddr:  "",
        SourcePort:  0,
        TargetAddr:  "127.0.0.1",
        TargetPort:  0,
        Protocol:    "tcp",
    }

    go func() {
        defer source.Close()
        _, err := io.WriteString(source, "Hello, world!")
        if err != nil {
            t.Error(err)
        }
    }()

    go func() {
        _, err := Forward(info)
        if err != nil {
            t.Errorf("forward failed: %v", err)
        }
    }()

    buffer := make([]byte, 1024)
    n, err := target.Read(buffer)
    if err != nil {
        t.Error(err)
    }

    if string(buffer[:n]) != "Hello, world!" {
        t.Errorf("Received incorrect data: %s", string(buffer[:n]))
    }

    time.Sleep(time.Second)
}
