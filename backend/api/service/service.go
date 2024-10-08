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

type IpApiResponse struct {
	Country    string
	RegionName string
	Query      string
}

type ProxyInfo struct {
	Type       string
	ExternalIp string
	Country    string
	Region     string
}

func getIpInfoWithProxy(proxy string, proxyType string) (result *IpApiResponse) {
	proxyUrl, error := url.Parse(proxyType + "://" + proxy)

	if error != nil {
		log.Printf("Bad proxy (%v) format: %e", proxy, error)
		return nil
	}

	httpClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   timeout * time.Second,
	}

	response, error := httpClient.Get(ipApiUrl)

	if error != nil {
		log.Printf("Proxy (%v) failed: %e", proxyUrl.String(), error)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Proxy (%v) returns bad status %s", proxyUrl.String(), response.StatusCode)
		return nil
	}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&result)

	log.Print(result)

	return result
}

func mapResponse(response *IpApiResponse, proxyType string) *ProxyInfo {
	return &ProxyInfo{
		Type:       proxyType,
		ExternalIp: response.Query,
		Country:    response.Country,
		Region:     response.RegionName,
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
