package main

import (
	"fmt"

	"github.com/das08/kuRakutanBot-go/module"
)

func main() {
	fmt.Println("Hello world!")
	r := module.LoadRakutanDetail()
	s, _ := r.Marshal()
	fmt.Println(fmt.Sprintf("%s", s))
}
