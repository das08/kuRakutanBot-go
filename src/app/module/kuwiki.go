package module

import (
	"encoding/json"
	"fmt"
	"github.com/das08/kuRakutanBot-go/models/kuwiki"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type KUWikiStatus struct {
	Success bool
	Result  string
}

func GetKakomonURL(e *Environments, lectureName string) KUWikiStatus {
	kakomonURL := ""
	method := "GET"
	req, err := http.NewRequest(method, e.KuwikiEndpoint, nil)
	if err != nil {
		log.Fatalf("NewRequest err=%s", err.Error())
	}

	q := req.URL.Query()
	q.Add("name", lectureName)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", e.KuwikiAccessToken))

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Client.Do err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}

	response := kuwiki.KUWiki{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("json.Unmarshal err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}

	for _, result := range response.Results {
		if result.Name == lectureName {
			for _, exam := range result.ExamSet {
				kakomonURL = exam.DriveLink
			}
		}
	}
	log.Println("[KUWiki] Got kakomon URL")
	return KUWikiStatus{true, kakomonURL}
}
