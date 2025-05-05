package main

import (
	"fmt"

	"github.com/hagelstam/forge/internal/http"
)

func main() {
	s := http.NewServer()

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
