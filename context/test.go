package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup
    wg.Add(2)
    go ope1(ctx, &wg)
    go ope2(ctx, &wg)

    time.Sleep(3 * time.Second)
    cancel()
    select {
    case <-ctx.Done():
        fmt.Println("ctx done.")
        wg.Wait()
        fmt.Println("programm end.")
    }
}

func ope1(ctx context.Context, wg *sync.WaitGroup) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("ope1 done.")
            wg.Done()
            return
        default:
            time.Sleep(1 * time.Second)
            fmt.Println("ope1 running")
        }
    }
}

func ope2(ctx context.Context, wg *sync.WaitGroup) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("ope2 done.")
            wg.Done()
            return
        default:
            time.Sleep(1 * time.Second)
            fmt.Println("ope2 running")
        }
    }
}

