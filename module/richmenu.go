package module

import (
	"io/ioutil"
	"log"

	"github.com/das08/ApexStalker-go/models/richmenu"
)

func LoadRakutanDetail() richmenu.RakutanDetail {
	jsonFile, err := ioutil.ReadFile("./assets/richmenu/rakutan_detail.json")
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := richmenu.UnmarshalRakutanDetail(jsonFile)

	return jsonData
}
