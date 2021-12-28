package main

import (
	"fmt"

	"github.com/das08/ApexStalker-go/module"
)

func main() {
	fmt.Println("Hello world!")
	r := module.LoadRakutanDetail()
	fmt.Println(r)
}
