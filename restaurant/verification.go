package restaurant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dev-hyunsang/waiting-services/config"
)

type RequestBusinessNubmer struct {
	BusinessNubmer []string `json:"b_no"`
}

// {"request_cnt":1,"match_cnt":1,"status_code":"OK","data":[{"b_no":"","b_stt":"계속사업자","b_stt_cd":"01","tax_type":"부가가치세 일반과세자","tax_type_cd":"01","end_dt":"","utcc_yn":"N","tax_type_change_dt":"","invoice_apply_dt":""}]}
type ResponseBusinessNubmer struct {
	RequestCNT int    `json:"request_cnt"`
	MatchCNT   int    `json:"match_cnt"`
	StatusCode string `json:"status_code"`
	Data       []BusinessData
}

type BusinessData struct {
	BusinessNubmer  string `json:"b_no"`
	BusinessSTT     string `json:"b_stt"`
	BusinessSTTCD   string `json:"b_stt_cd"`
	TaxType         string `json:"tax_type"`
	TaxTypeCd       string `json:"tax_type_cd"`
	EndDt           string `json:"end_dt"`
	UtccYn          string `json:"utcc_yn"`
	TaxTypeChangeDt string `json:"tax_type_change_dt"`
	InvoiceApplyDt  string `json:"invoice_apply_dt"`
}

func CheackingBusinessNumber(number string) bool {
	url := "https://api.odcloud.kr/api/nts-businessman/v1/status?serviceKey=" + config.GetEnv("API_KEY")
	method := "POST"

	rbn := RequestBusinessNubmer{
		BusinessNubmer: []string{
			number,
		},
	}

	pbytes, err := json.Marshal(rbn)
	if err != nil {
		log.Println("[ERROR] CheackingBusinessNumber | Failed to JSON Marshal")
		log.Println(err)
		return false
	}

	buff := bytes.NewBuffer(pbytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	responseJSON := new(ResponseBusinessNubmer)
	err = json.Unmarshal([]byte(body), responseJSON)
	if err != nil {
		log.Println("[ERROR] CheackingBusinessNumber | Failed to JSON Unmarshal")
	}

	log.Println(responseJSON)
	if responseJSON.StatusCode == "OK" {
		return true
	}

	return false
}
