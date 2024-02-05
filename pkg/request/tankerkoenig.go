package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestTankerkoenig struct {
	ApiKey string `json:"api_key"`
}

type Prices struct {
	Status string  `json:"status"`
	E5     float32 `json:"e5"`
	E10    float32 `json:"e10"`
	Diesel float32 `json:"diesel"`
}

type PricesRespond struct {
	Ok      bool              `json:"ok"`
	License string            `json:"license"`
	Data    string            `json:"data"`
	Prices  map[string]Prices `json:"prices"`
}

func (r *RequestTankerkoenig) MakeRequest(ids []string) (*PricesRespond, string, error) {
	if len(ids) > 10 {
		return nil, "", fmt.Errorf("too many ids in request /// Max 10 ids are allowed")
	}
	var idsString string
	for index, value := range ids {
		if index == 0 {
			idsString = value
		} else {
			idsString = fmt.Sprintf("%s,%s", idsString, value)
		}
	}
	url := fmt.Sprintf("https://creativecommons.tankerkoenig.de/json/prices.php?ids=%s&apikey=%s", idsString, r.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	completBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	respondStruct := new(PricesRespond)
	err = json.NewDecoder(bytes.NewReader(completBody)).Decode(respondStruct)
	if err != nil {
		return nil, "", err
	}
	return respondStruct, string(completBody), nil
}
