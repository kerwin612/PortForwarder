package main

import (
    "os"
    "fmt"
    "flag"
    "errors"
    "strings"
    "strconv"
    "syscall"
    "os/signal"
    "github.com/kerwin612/port-forwarder/lib/core"
)

func parseSource(s string) (string, int, error) {
	if lastIndex := strings.LastIndex(s, ":"); lastIndex != -1 {
		strPart := s[:lastIndex]
		intPart := s[lastIndex+1:]

		intVal, err := strconv.Atoi(intPart)
		if err != nil {
			return "", 0, errors.New("invalid port number")
		}

		return strPart, intVal, nil
	}

	return "", 0, errors.New("invalid source format")
}

func parseTarget(s string) (string, string, int, error) {
	firstIndex := strings.Index(s, ":")
	lastIndex := strings.LastIndex(s, ":")

	if firstIndex == -1 || lastIndex == -1 || firstIndex == lastIndex {
		return "", "", 0, errors.New("invalid target format")
	}

	strPart1 := s[:firstIndex]
	strPart2 := s[firstIndex+1 : lastIndex]
	intPart := s[lastIndex+1:]

	intVal, err := strconv.Atoi(intPart)
	if err != nil {
		return "", "", 0, errors.New("invalid port number")
	}

	return strPart1, strPart2, intVal, nil
}

func main() {
    source := flag.String("s", "", "Source Addr: addr:port")
    target := flag.String("t", "", "Target Addr: protocol:addr:port")

    flag.Parse()

    source_addr, source_port, err := parseSource(*source)
    if err != nil {
        panic(err)
    }

    protocol, target_addr, target_port, err := parseTarget(*target)
    if err != nil {
        panic(err)
    }

    ln, err := core.Forward(core.Info{
        SourceAddr:  source_addr,
        SourcePort:  source_port,
        TargetAddr:  target_addr,
        TargetPort:  target_port,
        Protocol:    protocol,
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Forwarded: [[%s]] ==>> [[%s]]\n", *source, *target)

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    sig := <-sigs
    fmt.Printf("Received signal: %s, shutting down...\n", sig)

    ln.Close()
    os.Exit(0)
}
