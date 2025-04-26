package api

import (
	"backend/config"
	"backend/models"
	"backend/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResponeStruct struct {
	Items    []StockApi `json:"items"`
	NextPage string     `json:"next_page"`
}

type StockApi struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

func FetchData() (string, error) {
	baseURl := config.LoadApi().URL
	token := config.LoadApi().Token
	var batch []models.Stock
	nextPage := ""
	resp := ""
	for {
		url := baseURl
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
		}
		items, newNextPage, err := FetchPage(url, token)
		if err != nil {
			return "", fmt.Errorf("Can't get page: Error %v", err)
		}

		for _, item := range items {
			//fmt.Println("this is the data :", item)
			stock, err := ConvertStockApi(item)
			if err == nil {
				batch = append(batch, stock)
			}
		}


		for len(batch) >= 100 {
			resp, err = repositories.StoreStock(batch)
			batch = batch[len(batch):]
		}

		if newNextPage == "" {
			for len(batch) > 0 {
				repositories.StoreStock(batch)
				batch = batch[len(batch):]
				fmt.Println("this is the final batch :", batch)
			}
			break
		}
		fmt.Println("this is the next page :", newNextPage)
		nextPage = newNextPage
	}

	return resp, nil
}

func FetchPage(url, token string) ([]StockApi, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("fail Request: %v", err)
	}

	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("fail doing request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("fail reading the response: %v", err)
	}

	var data ResponeStruct
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Errorf("fail reading the body: %v", err)
	}
	return data.Items, data.NextPage, nil
}

func ConvertStockApi(stocktoConvert StockApi) (models.Stock, error) {

	targetFromStr := stocktoConvert.TargetFrom
	targetToStr := stocktoConvert.TargetTo

	targetFromStr = strings.TrimPrefix(targetFromStr, "$")
	targetFromStr = strings.ReplaceAll(targetFromStr, ",", "")

	targetToStr = strings.TrimPrefix(targetToStr, "$")
	targetToStr = strings.ReplaceAll(targetToStr, ",", "")

	targetFrom, err1 := strconv.ParseFloat(targetFromStr, 64)
	targetTo, err2 := strconv.ParseFloat(targetToStr, 64)
	time, err3 := time.Parse(time.RFC3339, stocktoConvert.Time)

	if err1 != nil || err2 != nil || err3 != nil {
		return models.Stock{}, fmt.Errorf("invalid data in API response %v %v %v", err1, err2, err3)
	}

	return models.Stock{
		Ticker:     stocktoConvert.Ticker,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Company:    stocktoConvert.Company,
		Action:     stocktoConvert.Action,
		Brokerage:  stocktoConvert.Brokerage,
		RatingFrom: stocktoConvert.RatingFrom,
		RatingTo:   stocktoConvert.RatingTo,
		Time:       time,
	}, nil
}
