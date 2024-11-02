package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"
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

func getIpInfoWithProxy(proxy string, proxyType string, ctx context.Context) *IpApiResponse {
	proxyUrl, error := url.Parse(proxyType + "://" + proxy)

	if error != nil {
		log.Printf("Bad proxy (%s) format: %e", proxy, error)
		return nil
	}

	httpClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   timeout * time.Second,
	}

	request, error := http.NewRequestWithContext(ctx, "GET", ipApiUrl, nil)
	if error != nil {
		log.Printf("Bad proxy (%s) format: %e", proxy, error)
		return nil
	}

	start := time.Now()
	response, error := httpClient.Do(request)
	past := time.Since(start).Seconds()

	if error != nil {
		log.Printf("Proxy (%s) failed: %e", proxyUrl.String(), error)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Proxy (%s) returns bad status %d", proxyUrl.String(), response.StatusCode)
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
	proxyTypes := [3]string{"sock5", "https", "http"}

	proxyInfoChannel := make(chan *ProxyInfo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	for _, proxyType := range proxyTypes {
		wg.Add(1)

		go func() {
			defer wg.Done()

			response := getIpInfoWithProxy(proxy, proxyType, ctx)

			if response != nil {
				proxyInfoChannel <- mapResponse(response, proxyType)
			}
		}()
	}

	go func() {
		wg.Wait()
		proxyInfoChannel <- nil
		close(proxyInfoChannel)
	}()

	return <-proxyInfoChannel
}
