package module

import (
	"encoding/json"
	"fmt"
	"github.com/das08/kuRakutanBot-go/models/kuwiki"
	"io/ioutil"
	"log"
	"net/http"
)

func GetKakomonURL(e *Environments, lectureName string) *string {
	var kakomonURL *string = nil
	method := "GET"
	req, err := http.NewRequest(method, e.KUWIKI_ENDPOINT, nil)
	if err != nil {
		log.Fatalf("NewRequest err=%s", err.Error())
	}

	q := req.URL.Query()
	q.Add("name", lectureName)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", e.KUWIKI_ACCESS_TOKEN))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Client.Do err=%s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}

	response := kuwiki.KUWiki{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("json.Unmarshal err=%s", err.Error())
	}

	for _, result := range response.Results {
		if result.Name == lectureName {
			for _, exam := range result.ExamSet {
				kakomonURL = &exam.DriveLink
			}
		}
	}

	return kakomonURL
}
