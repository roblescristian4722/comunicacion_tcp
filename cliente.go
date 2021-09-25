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
    res := proceso.Proc{ 0, 0, false }
    c, err := net.Dial("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }
    // Canal que detecta el Ctrl + C
    sig := make(chan os.Signal, 1)
    killC := make(chan uint64)
    retC := make(chan uint64)
    signal.Notify(sig, os.Interrupt)

    gob.NewEncoder(c).Encode(res)
    gob.NewDecoder(c).Decode(&res)

    go proceso.Proceso(res.Id, res.Data, killC, retC)

    c, err = net.Dial("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }

    <-sig
    killC <- res.Id
    res.Data, _ = <-retC
    gob.NewEncoder(c).Encode(res)
    c.Close()
    os.Exit(1)
}

func main() {
    go cliente()
    fmt.Scanln()
}
