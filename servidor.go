package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "strings"
    "strconv"
    "./proceso"
)

type ProcStruct [2]uint64

func servidor() {
    var msg string
    p := make([]ProcStruct, 5)
    retC := make(chan uint64)
    killC := make(chan uint64)
    s, err := net.Listen("tcp", ":9999")
    defer s.Close()

    if err != nil {
        fmt.Println(err)
        return
    }
    
    // Inicializa goroutines
    for i, _ := range p { go proceso.Proceso(uint64(i), 0, killC, retC) }
    // Escucha peticiones de clientes
    for {
        c, err := s.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        gob.NewDecoder(c).Decode(&msg)
        fmt.Println(msg)
        res := strings.Split(msg, "|")
        if res[2] == "0" {
            go startClient(c, p, retC, killC)
        } else { go endClient(c, p, retC, killC, res) }
    }
}

func endClient(c net.Conn, p []ProcStruct, retC chan uint64, killC chan uint64, res []string) {
    c.Close()
    id, _ := strconv.ParseUint(res[0], 10, 64)
    data, _ := strconv.ParseUint(res[1], 10, 64)
    p[id] = [2]uint64{ 0, 0 }
    go proceso.Proceso(id, data, killC, retC)
}

func startClient(c net.Conn, p []ProcStruct, retC chan uint64, killC chan uint64) {
    for i := 0; i < len(p); i++ {
        if p[i][1] == 0 {
            fmt.Println("i: ", i)
            killC <- uint64(i)
            data := <-retC
            p[i][1] = 1
            gob.NewEncoder(c).Encode(strconv.FormatUint(uint64(i), 10) + "|" + strconv.FormatUint(data, 10) + "|" + "1")
            break
        }
    }
    c.Close()
}

func main() {
    go servidor()
    fmt.Scanln()
}
