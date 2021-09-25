package proceso

import (
    "fmt"
    "time"
)

var endl int = 0
var max int = 0

func Proceso(id uint64, i uint64, kill chan uint64, ret chan uint64) {
    max ++
	for {
        select {
            case res := <-kill:
                if res == id {
                    ret <- i
                    max--
                    return
                } else { kill <- res }
            case <-time.After(time.Millisecond * 500):
                fmt.Printf("\nid %d: %d", id, i)
                endl = ( endl + 1 ) % max
                if endl == 0 && max > 1 { fmt.Print("\n----------") }
                i = i + 1
        }
    }
}
