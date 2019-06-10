package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/gohargasparyan/access-please"
	"github.com/gohargasparyan/access-please/common"
)

func main() {
	var context = kingpin.Flag("context", "Defines kubernetes context.").Required().PlaceHolder("<context>").Short('c').String()

	kingpin.Parse()
	fmt.Printf("%s\n", *context)

	ap, err := accessplease.New(*context)
	common.Panic(err)

	ap.Run()
}
