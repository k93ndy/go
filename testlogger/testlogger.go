package main

import (
    "strconv"
//    "fmt"
    "os"
    "time"
    "log"
    "io"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    testlog, err := os.Create("/tmp/testlog")
    check(err)
    defer testlog.Close()

    errorlog, err := os.Create("/tmp/errorlog")
    check(err)
    defer errorlog.Close()

    Info := log.New(io.MultiWriter(os.Stdout, testlog), "Info:", log.Ldate | log.Ltime| log.Lshortfile)
    Error := log.New(io.MultiWriter(os.Stderr, errorlog), "Error:", log.Ldate | log.Ltime| log.Lshortfile)

    iter := 1
    for ;; {
        msg := "Count: " + strconv.Itoa(iter)
        Info.Println(msg)
        if iter%10 == 0 {
            msg := strconv.Itoa(iter) + " is a multiple of 10"
            Error.Println(msg)
        }
        iter += 1
        time.Sleep(3 * time.Second)
    }
}
