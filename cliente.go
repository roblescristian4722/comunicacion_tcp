package main
import (
    "fmt"
    "net"
    "encoding/gob"
    "os"
    "os/signal"
    "./proceso"
)

func cliente() {
    kill := false
    res := proceso.Proc{ 0, 0, false }
    c, err := net.Dial("tcp", ":9999")
    defer c.Close()

    if err != nil {
        fmt.Println(err)
        return
    }

    // Canal que detecta el Ctrl + C
    sig := make(chan os.Signal, 1)
    killC := make(chan uint64)
    retC := make(chan uint64)
    signal.Notify(sig, os.Interrupt)
    go func () {
        <-sig
        kill = true
        os.Exit(1)
    }()

    gob.NewEncoder(c).Encode(res)
    gob.NewDecoder(c).Decode(&res)
    fmt.Println(res)
    go proceso.Proceso(res.Id, res.Data, killC, retC)
    for {
        if kill {
            killC <- res.Id
            res.Data, _ = <-retC
            gob.NewEncoder(c).Encode(res)
            break
        }
    }
}

func main() {
    go cliente()
    fmt.Scanln()
}
