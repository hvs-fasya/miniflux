package alert

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/storage"
	"github.com/miniflux/miniflux/timer"
)

var (
	errRequestFailed    = "Unable to execute request: %v"
	errServerFailure    = "Unable to fetch feed (statusCode=%d)"
	errDuplicate        = "This feed already exists (%s)"
	errNotFound         = "Feed %d not found"
	errEncoding         = "Unable to normalize encoding: %q"
	errCategoryNotFound = "Category not found for this user"
	errEmptyFeed        = "This feed is empty"
	ADVISORIES          = "https://travel.gc.ca/destinations/"
)

// Handler contains all the logic to create and refresh alert.
type Handler struct {
	store *storage.Storage
}

func Task(store *storage.Storage) error {
	countries, err := store.Countries()
	handler := new(Handler)
	handler.store = store
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(countries))
	_, err = handler.CreateAlert(12)
	if err != nil {
		logger.Debug("error: %s", err)
		return err
	}
	return nil
}

func (h *Handler) CreateAlert(countryID int64) (*model.Alert, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Handler:CreateAlert] countryID=%s", countryID))

	alert := new(model.Alert)
	alert.CountryID = countryID

	country_name := "India"
	//country_name := "SURINAME"

	doc, err := goquery.NewDocument(ADVISORIES + strings.ToLower(country_name))
	if err != nil {
		logger.Debug("error: %s", err)
		return alert, nil
	}

	alert.LastUpdated, err = time.Parse("January 02, 2006 15:04", doc.Find("span#Label9").Text())
	if err != nil {
		logger.Debug("error: %s", err)
		return alert, nil
	}
	alert.StillValid, err = time.Parse("January 02, 2006 15:04", doc.Find("span#Label12").Text())
	if err != nil {
		logger.Debug("error: %s", err)
		return alert, nil
	}
	alert.LatestUpdates = doc.Find("span#Label11").Text()
	doc.Find("div.AdvisoryContainer").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			alert.RiskLevel = s.Find("h3").Text()
			alert.RiskDetails = s.Find("p").Text()
		}
	})
	healthTab := doc.Find("#health")
	panel := healthTab.Find(".panel-body")
	panel.Find("li").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			healthLink := s.Find("a")
			alert.HealthTitle = strings.Replace(healthLink.Text(), ": Advice for travellers", "", -1)
			alert.HealthDate, err = time.Parse("January 02, 2006", strings.Replace(s.Text(), healthLink.Text()+" - ", "", -1))
			if err != nil {
				logger.Debug("error: %s", err)
			}
			contentLink, _ := healthLink.Attr("href")
			healthDoc, err := goquery.NewDocument(contentLink)
			if err != nil {
				logger.Debug("error: %s", err)
			}
			alert.HealthContent, err = healthDoc.Find("main").Html()
			if err != nil {
				logger.Debug("error: %s", err)
			}
		}
	})
	if err != nil {
		logger.Debug("error: %s", err)
	}

	err = h.store.CreateAlert(alert)
	if err != nil {
		return nil, err
	}

	logger.Debug("[Handler:CreateAlert] Alert saved for countryID: %d", alert.CountryID)

	return alert, nil
}

//doc, err := goquery.NewDocument("https://travel.gc.ca/travelling/advisories")
//if err != nil {
//logger.Debug("error: %s", err)
//html.ServerError(w, err)
//return
//}
//var data = make(map[string]OutputCountry)
//
//var isos []ISO
//file, err := ioutil.ReadFile("./news/3-codes-news.json")
//if err != nil {
//fmt.Println("error: %s", err)
//}
//err = json.Unmarshal(file, &isos)
//isomap := make(map[string]string)
//for _, el := range isos {
//isomap[el.Name] = el.Alpha3
//}
//
//trs := doc.Find("tr.gradeX")
//trs.Each(func(i int, tr *goquery.Selection) {
//tds := tr.Find("td")
//var countryName string
//var countryCode string
//tds.Each(func(j int, td *goquery.Selection) {
//if j == 1 {
//countryName = td.Text()
//_, ok := isomap[countryName]
//if !ok {
//fmt.Println("ERROR", countryName)
//} else {
//countryCode, _ = isomap[countryName]
//}
//}
//if j == 2 {
//riskLevel, _ := levels[td.Text()]
//data[countryCode] = OutputCountry{riskLevel, td.Text()}
//}
//})
//
//})
