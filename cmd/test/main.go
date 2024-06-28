package main

import (
	"fmt"
	"github.com/lgrisa/lib/utils"
)

func main() {
	fmt.Printf("%s@%s", utils.GetUsername(), utils.GetHostname())
}
