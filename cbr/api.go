package cbr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Daily(date time.Time) (*DailyResponse, error) {
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date.Format("02/01/2006"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var daily DailyResponse
	err = xmlDecoder(body).Decode(&daily)
	if err != nil {
		return nil, err
	}
	return &daily, nil
}
