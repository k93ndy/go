package server

import (
    "context"
    "log"
    "net"
    "sync"
)

type Conn struct {
    svr         *Server
    ctx         *contexts
    conn        *net.TCPConn
    readCtx     context.Context
    stopRead    context.CancelFunc
    writeCtx    context.Context
    stopWrite   context.CancelFunc
    sem         chan struct{}
    wg          sync.WaitGroup
}

func newConn(svr *Server, ctx *contexts, tcpConn *net.TCPConn) *Conn {
    readCtx, stopRead := context.WithCancel(context.Background())
    writeCtx, stopWrite := context.WithCancel(context.Background())
    sem := make(chan struct{}, 1)
    return &Conn{
        svr:     svr,
        ctx:     ctx,
        conn:    tcpConn,
        readCtx: readCtx,
        stopRead: stopRead,
        writeCtx: writeCtx,
        stopWrite: stopWrite,
        sem: sem,
    }
}

func (c *Conn) handleConnection() {
    defer func() {
        c.stopWrite()
        delete(c.svr.Conns, c.conn.RemoteAddr())
        // for clientAddr,_ := range c.svr.Conns{
        //     log.Printf("client: %v", clientAddr)
        // }
        c.conn.Close()
        log.Printf("%v: disconnected", c.conn.RemoteAddr())
        c.svr.Wg.Done()
    }()

    log.Printf("%v: connected", c.conn.RemoteAddr())
    go c.handleRead()

    select {
    case <-c.readCtx.Done():
    case <-c.ctx.ctxShutdown.Done():
    case <-c.svr.AcceptCtx.Done():
    case <-c.ctx.ctxGraceful.Done():
        c.conn.CloseRead()
        c.wg.Wait()
    }
}

func (c *Conn) handleRead() {
    defer c.stopRead()

    buf := make([]byte, 4*1024)

    for {
        n, err := c.conn.Read(buf)
        if err != nil {
            if ne, ok := err.(net.Error); ok {
                switch {
                case ne.Temporary():
                    continue
                }
            }
            log.Printf("%v: Read %v", c.conn.RemoteAddr(), err)
            return
        }

        wBuf := make([]byte, n)
        copy(wBuf, buf[:n])
        c.wg.Add(1)
        // go c.handleEcho(wBuf)
        go c.broadcast(wBuf)
    }
}

func (c *Conn) handleEcho(buf []byte) {
    // defer c.wg.Done()

    for {
        select {
        case <-c.writeCtx.Done():
            return
        case c.sem <- struct{}{}:
            defer func() { <-c.sem }()
            n, err := c.conn.Write(buf)
            if err != nil {
                if nerr, ok := err.(net.Error); ok {
                    if nerr.Temporary() {
                        buf = buf[n:]
                        continue
                    }
                }
                log.Printf("%v: Write error %v", c.conn.RemoteAddr(), err)
                // write error
                c.stopRead()
                c.stopWrite()
            }
            return
        }
    }
}

func (c *Conn) broadcast(buf []byte) {
    defer func() {
        c.wg.Done()
        // log.Printf("broadcast done.")
    }()

    for _,conn := range c.svr.Conns {
        go conn.handleEcho(buf)
    }
}
