package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "./proceso"
)

type ProcStruct [2]uint64

func servidor() {
    p := make([]ProcStruct, 5)
    retC := make(chan uint64)
    killC := make(chan uint64)
    res := proceso.Proc{}
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
        gob.NewDecoder(c).Decode(res)

        if err != nil {
            fmt.Println(err)
            continue
        }
        if res.Active == true {
            go endClient(c, p, retC, killC, res)
        } else {
            fmt.Println("res: ", res)
            go startClient(c, p, retC, killC)
        }
    }
}

func endClient(c net.Conn, p []ProcStruct, retC chan uint64, killC chan uint64, res proceso.Proc) {
    
}

func startClient(c net.Conn, p []ProcStruct, retC chan uint64, killC chan uint64) {
    for i := 0; i < len(p); i++ {
        if p[i][1] == 0 {
            killC <- uint64(i)
            data := <-retC
            p[i][1] = 1
            fmt.Println([3]uint64{uint64(i), data, 1})
            gob.NewEncoder(c).Encode(proceso.Proc{ 
                                                Id: uint64(i),
                                                Data: data,
                                                Active: true })
            break
        }
    }
}

func main() {
    go servidor()
    fmt.Scanln()
}
