package main
import (
    "fmt"
    "net"
    "encoding/gob"
    "os"
    "os/signal"
    "strings"
    "strconv"
    "./proceso"
)

func cliente() {
    msg := "0|0|0"
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

    enc := gob.NewEncoder(c)
    dec := gob.NewDecoder(c)

    enc.Encode(msg)
    dec.Decode(&msg)
    c.Close()

    res := strings.Split(msg, "|")
    id, _ := strconv.ParseUint(res[0], 10, 64)
    data, _ := strconv.ParseUint(res[1], 10, 64)

    go proceso.Proceso(id, data, killC, retC)

    // Bloque que se ejecutará cuando el usuario presione Ctrl + C
    <-sig
    killC <- id
    data, _ = <-retC

    c, err = net.Dial("tcp", ":9999")
    msg2 := strconv.FormatUint(id, 10) + "|" + strconv.FormatUint(data, 10) + "|" + "1"
    gob.NewEncoder(c).Encode(msg2)
    c.Close()
    os.Exit(1)
}

func main() {
    go cliente()
    fmt.Scanln()
}
