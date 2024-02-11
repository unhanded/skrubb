package main

import (
	"fmt"

	"github.com/unhanded/skrubb/internal/app/libskrubb"
)

func main() {
	c := libskrubb.AttatchToContainer("S1")
	fmt.Printf("%s", c.Name)
}
