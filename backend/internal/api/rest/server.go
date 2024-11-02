package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"proxy-checker-server/internal/service"
	"sync"
)

type ChechRequest struct {
	Value []string `json:"value"`
}

type ChechResponse struct {
	CheckResult []*service.ProxyInfo `json:"checkResult"`
}

func Check(w http.ResponseWriter, response *http.Request) {
	defer response.Body.Close()

	var data = ChechRequest{}
	json.NewDecoder(response.Body).Decode(&data)

	wg := sync.WaitGroup{}

	var result = make([]*service.ProxyInfo, 0, len(data.Value))

	for _, value := range data.Value {
		log.Printf("Handling proxy: %v", value)

		wg.Add(1)

		go func() {
			log.Printf("Handling proxy: %v", value)

			defer wg.Done()

			proxyInfo := service.GetProxyInfo(value)

			if proxyInfo != nil {
				result = append(result, proxyInfo)
			}
		}()
	}

	wg.Wait()

	jsonResp, err := json.Marshal(ChechResponse{CheckResult: result})

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func StartServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, response *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", response.URL.Path)
	})

	http.HandleFunc("/check", Check)

	log.Printf("Listening http on %s", port)
	http.ListenAndServe(":"+port, nil)
}
