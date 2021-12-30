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
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}

func LoadSearchResult() richmenu.SearchResult {
	jsonFile, err := ioutil.ReadFile("./assets/richmenu/search_result.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := richmenu.UnmarshalSearchResult(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}

func LoadSearchResultMore() richmenu.SearchResultMore {
	jsonFile, err := ioutil.ReadFile("./assets/richmenu/search_result_more.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := richmenu.UnmarshalSearchResultMore(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}

func LoadFavorites() richmenu.Favorites {
	jsonFile, err := ioutil.ReadFile("./assets/richmenu/favorites.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := richmenu.UnmarshalFavorites(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}

func LoadHelp() richmenu.Help {
	jsonFile, err := ioutil.ReadFile("./assets/richmenu/help.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := richmenu.UnmarshalHelp(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}
