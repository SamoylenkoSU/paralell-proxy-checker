package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	ipApiUrl = "http://ip-api.com/json"
	timeout  = 30
)

type IpApiData struct {
	Country string
	City    string
	Query   string
}

type IpApiResponse struct {
	Data    IpApiData
	Timeout float64
}

type ProxyInfo struct {
	Type       string
	ExternalIp string
	Country    string
	City       string
	Timeout    float64
}

func getIpInfoWithProxy(proxy string, proxyType string) *IpApiResponse {
	proxyUrl, error := url.Parse(proxyType + "://" + proxy)

	if error != nil {
		log.Printf("Bad proxy (%v) format: %e", proxy, error)
		return nil
	}

	httpClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   timeout * time.Second,
	}

	start := time.Now()
	response, error := httpClient.Get(ipApiUrl)
	past := time.Since(start).Seconds()

	if error != nil {
		log.Printf("Proxy (%v) failed: %e", proxyUrl.String(), error)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Proxy (%v) returns bad status %s", proxyUrl.String(), response.StatusCode)
		return nil
	}

	defer response.Body.Close()

	data := IpApiData{}

	json.NewDecoder(response.Body).Decode(&data)

	return &IpApiResponse{
		Data:    data,
		Timeout: past,
	}
}

func mapResponse(response *IpApiResponse, proxyType string) *ProxyInfo {
	return &ProxyInfo{
		Type:       proxyType,
		ExternalIp: response.Data.Query,
		Country:    response.Data.Country,
		City:       response.Data.City,
		Timeout:    response.Timeout,
	}
}

func GetProxyInfo(proxy string) *ProxyInfo {
	var response *IpApiResponse

	proxyTypes := [3]string{"sock5", "https", "http"}

	proxyInfoChannel := make(chan *ProxyInfo)

	counter := 0
	isHandled := false

	for _, proxyType := range proxyTypes {
		go func() {
			response = getIpInfoWithProxy(proxy, proxyType)

			if response != nil {
				if isHandled {
					return
				}
				proxyInfoChannel <- mapResponse(response, proxyType)
				isHandled = true
				close(proxyInfoChannel)
				return
			}

			counter++

			if counter >= len(proxyTypes) {
				proxyInfoChannel <- nil
				close(proxyInfoChannel)
			}
		}()
	}

	return <-proxyInfoChannel
}
