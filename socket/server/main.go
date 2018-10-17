package main

import (
    "log"
    "flag"
    "os"
    "os/signal"
    "context"
    "syscall"
    "./server"
)

type ServerOpts struct {
    listenAddr string
}

func (o *ServerOpts) loadConf() {
    // TODO: load config from file or env
    o.listenAddr = "localhost:30000"
}

func (o *ServerOpts) setFlags(f *ServerOpts) {
    o.listenAddr = f.listenAddr
}

func main() {
    serverOpts := &ServerOpts{}
    serverOpts.loadConf()
    listenAddr := flag.String("l", "localhost:30000", "listen address:port")
    flag.Parse()
    serverOpts.setFlags(&ServerOpts{listenAddr: *listenAddr})

    sigChan := make(chan os.Signal, 1)
    // Ignore all signals
    signal.Ignore()
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

    //host := loadConf(*listenAddr)
    svr := server.NewServer(context.Background(), serverOpts.listenAddr)
    err := svr.Listen()
    if err != nil {
        log.Fatalln(err)
    }

    log.Println("Server Started")

    for {
        select {
        case sig := <-sigChan:
            switch sig {
            case syscall.SIGINT, syscall.SIGTERM:
                log.Println("Server Shutdown...")
                svr.Shutdown()
                svr.Wg.Wait()
                <-svr.ChClosed
                log.Println("Server Shutdown Completed")
            case syscall.SIGQUIT:
                log.Println("Server Graceful Shutdown...")
                svr.GracefulShutdown()
                svr.Wg.Wait()
                <-svr.ChClosed
                log.Println("Server Graceful Shutdown Completed")
            case syscall.SIGHUP:
                log.Println("Server Restarting...")

                serverOpts.loadConf()

                svr, err = svr.Restart(context.Background(), serverOpts.listenAddr)
                if err != nil {
                    log.Fatal(err)
                }
                log.Println("Server Restarted")
                continue
            default:
                panic("unexpected signal has been received")
            }
        case <-svr.AcceptCtx.Done():
            log.Println("Server Error Occurred")
            // wait until all connection closed
            svr.Wg.Wait()
            <-svr.ChClosed
            log.Println("Server Shutdown Completed")
        }
        return
    }
}
