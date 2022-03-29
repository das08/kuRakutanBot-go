package module

import (
	"io/ioutil"
	"log"

	"github.com/das08/kuRakutanBot-go/models/richmenu"
)

func LoadRakutanDetail() richmenu.RakutanDetail {
	jsonFile := LoadJSON("./assets/richmenu/rakutan_detail.json")
	jsonData, err := richmenu.UnmarshalRakutanDetail(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadSearchResult() richmenu.SearchResult {
	jsonFile := LoadJSON("./assets/richmenu/search_result.json")
	jsonData, err := richmenu.UnmarshalSearchResult(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadSearchResultMore() richmenu.SearchResultMore {
	jsonFile := LoadJSON("./assets/richmenu/search_result_more.json")
	jsonData, err := richmenu.UnmarshalSearchResultMore(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadFavorites() richmenu.Favorites {
	jsonFile := LoadJSON("./assets/richmenu/favorites.json")
	jsonData, err := richmenu.UnmarshalFavorites(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadJSON(filename string) []byte {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return jsonFile
}
