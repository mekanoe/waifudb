package main

import (
	"fmt"
	"log"

	"github.com/kayteh/waifudb/cmd/waifudb/run"
	"github.com/kayteh/waifudb/datastore"
)

func main() {
	_, err := datastore.New(nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("こんばんわ")

	run.Start("localhost:7099")
}
