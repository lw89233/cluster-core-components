package recovery

import (
	"log"
)

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC_RECOVERED\t%v\n", err)
			}
		}()
		fn()
	}()
}