package main

import (
	"fmt"
	"os"

	"github.com/unhanded/skrubb/internal/app/libskrubb/container"
)

func main() {
	err := container.MakeContainer()
	if err != nil {
		fmt.Printf("ERROR - %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
