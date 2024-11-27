package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/emicklei/xml2go"
)

var xmlDoc = flag.String("xml", "", "path to xml file")

func main() {
	flag.Parse()
	if *xmlDoc == "" {
		log.Println("missing -xml option")
		return
	}
	f, _ := os.Open(*xmlDoc)
	b := xml2go.NewBuilder()
	if err := b.parse(f); err != nil {
		slog.Error("fail", "err", err.Error())
	}
	for _, each := range b.StructsMap {
		fmt.Println(each.String())
	}
}
