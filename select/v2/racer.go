package racer

import (
	"fmt"
	"net/http"
)

// Racer compares the response times of a and b, returning the fastest one.
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		fmt.Println("returning a")
		return a
	case <-ping(b):
		fmt.Println("returning b")
		return b
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
