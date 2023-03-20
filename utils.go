package main

import (
	"log"
	"time"
)

func timer[V any](name string, caller func() V) V {
	start := time.Now()
	callback := caller()
	elapsedTime := time.Since(start)

	enableLog := true
	if enableLog {
		log.Printf("%s took %s", name, elapsedTime)
	}
	return callback
}
