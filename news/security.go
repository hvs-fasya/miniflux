package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
)

var (
	levels = map[string]string{
		"Exercise normal security precautions":                            "NormalPrecautions",
		"Exercise normal security precautions (with regional advisories)": "NormalPrecautionsReg",
		"Exercise a high degree of caution":                               "HighCaution",
		"Exercise a high degree of caution (with regional advisories)":    "HighCautionReg",
		"Avoid non-essential travel":                                      "AvoidNonEssential",
		"Avoid non-essential travel (with regional advisories)":           "AvoidNonEssentialReg",
		"Avoid all travel (with regional advisories)":                     "AvoidAllReg",
		"Avoid all travel":                                                "AvoidAll",
	}
)

type OutputCountry struct {
	FillKey string `json:"fillKey"`
	Risk    string `json:"risk"`
}

type ISO struct {
	Name        string `json:"name"`
	Alpha3      string `json:"alpha-3"`
	CountryCode string `json:"country-code"`
}

// Security scraps and  sends security data
func (c *Controller) Security(w http.ResponseWriter, r *http.Request) {
	doc, err := goquery.NewDocument("https://travel.gc.ca/travelling/advisories")
	if err != nil {
		logger.Debug("error: %s", err)
		html.ServerError(w, err)
		return
	}
	var data = make(map[string]OutputCountry)

	var isos []ISO
	file, err := ioutil.ReadFile("./news/3-codes-news.json")
	err = json.Unmarshal(file, &isos)
	isomap := make(map[string]string)
	for _, el := range isos {
		isomap[el.Name] = el.Alpha3
	}

	trs := doc.Find("tr.gradeX")
	trs.Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find("td")
		var countryName string
		var countryCode string
		tds.Each(func(j int, td *goquery.Selection) {
			if j == 1 {
				countryName = td.Text()
				_, ok := isomap[countryName]
				if !ok {
					fmt.Println("ERROR", countryName)
				} else {
					countryCode, _ = isomap[countryName]
				}
			}
			if j == 2 {
				riskLevel, _ := levels[td.Text()]
				data[countryCode] = OutputCountry{riskLevel, td.Text()}
			}
		})

	})
	js, _ := json.Marshal(data)

	//dataStr := `{
	//"USA": { "fillKey": "AvoidAll" },
	//"JPN": { "fillKey": "NormalPrecautionsReg" },
	//"CAN": { "fillKey": "NormalPrecautions" },
	//"RUS": { "fillKey": "HighCaution" },
	//"IND": { "fillKey": "HighCautionReg" }
	//}`

	w.Write(js)
	return
}
