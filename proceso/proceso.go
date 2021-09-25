package proceso

import (
    "fmt"
    "time"
)

type Proc struct {
    Id uint64
    Data uint64
    Active bool
}

func Proceso(id uint64, i uint64, kill chan uint64, ret chan uint64) {
	for {
        select {
            case res := <-kill:
                if res == id {
                    ret <- i
                    return
                } else { kill <- res }
            default:
                fmt.Printf("\nid %d: %d", id, i)
                i = i + 1
                time.Sleep(time.Millisecond * 500)
        }
    }
}