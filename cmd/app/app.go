package main

import (
	"fmt"

	"github.com/das08/kuRakutanBot-go/module"
)

func main() {
	env := module.LoadEnv(true)
	mongo := module.CreateDBClient(&env)
	defer mongo.Cancel()
	defer mongo.Client.Disconnect(mongo.Ctx)
	fmt.Println("Hello world!")
	// r := module.LoadRakutanDetail()
	// s, _ := r.Marshal()
	// fmt.Println(fmt.Sprintf("%s", s))
	// module.CreateDBClient(&env)

	module.FindOne(&env, mongo)
}
