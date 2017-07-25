package main

import (
	"fmt"
	"log"

	"github.com/kayteh/waifudb/cmd/waifudb/run"
)

func main() {
	fmt.Println("こんばんわ")

	s, err := run.New(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	s.ListenAndServe()
}
