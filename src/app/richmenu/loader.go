package richmenu

import (
	"io/ioutil"
	"log"
)

var (
	RakutanDetailJson    RakutanDetail
	SearchResultJson     SearchResult
	SearchResultMoreJson SearchResultMore
	FavoritesJson        Favorites
)

func PreloadJson() {
	var jsonFile []byte
	var err error
	jsonFile = LoadJSON("./assets/richmenu/rakutan_detail.json")
	RakutanDetailJson, err = UnmarshalRakutanDetail(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonFile = LoadJSON("./assets/richmenu/search_result.json")
	SearchResultJson, err = UnmarshalSearchResult(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonFile = LoadJSON("./assets/richmenu/search_result_more.json")
	SearchResultMoreJson, err = UnmarshalSearchResultMore(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonFile = LoadJSON("./assets/richmenu/favorites.json")
	FavoritesJson, err = UnmarshalFavorites(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadJSON(filename string) []byte {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return jsonFile
}
