package richmenu

import (
	"io/ioutil"
	"log"
)

func LoadRakutanDetail() RakutanDetail {
	jsonFile := LoadJSON("./assets/richmenu/rakutan_detail.json")
	jsonData, err := UnmarshalRakutanDetail(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadSearchResult() SearchResult {
	jsonFile := LoadJSON("./assets/richmenu/search_result.json")
	jsonData, err := UnmarshalSearchResult(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadSearchResultMore() SearchResultMore {
	jsonFile := LoadJSON("./assets/richmenu/search_result_more.json")
	jsonData, err := UnmarshalSearchResultMore(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func LoadFavorites() Favorites {
	jsonFile := LoadJSON("./assets/richmenu/favorites.json")
	jsonData, err := UnmarshalFavorites(jsonFile)
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
