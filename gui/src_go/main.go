package src_go

import (
    "io"
    "os"
    "log"
    "flag"
    "os/user"
    "path/filepath"
    "runtime/debug"
    _ "github.com/kerwin612/port-forwarder/gui/_statik"
)

var logger *log.Logger

func Main() {

    myself, err := user.Current()
    if err != nil {
        panic(err)
    }

    autoopen := flag.Bool("ao", true, "Does the program auto-open on launch")
    workspace := flag.String("ws", filepath.Join(myself.HomeDir, ".pf"), "The workspace where the program runs")

    flag.Parse()

    homedir, err := filepath.Abs(*workspace)
    if err != nil {
        panic(err)
    }

    if err := os.MkdirAll(homedir, 0775); err != nil {
        panic(err)
    }

    lf, err := os.OpenFile(filepath.Join(homedir, "log"), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
    if err != nil {
        panic(err)
    }

    logger = log.New(io.MultiWriter(lf, os.Stdout), "", log.LstdFlags)

    defer func() {
        if r := recover(); r != nil {
            logger.Printf("caught panic: %v\n", r)
            logger.Printf("stack trace:\n%s", debug.Stack())
        }
    }()

    logger.Printf("workspace: [%s].\n", homedir)

    Startup(homedir, *autoopen)

}
