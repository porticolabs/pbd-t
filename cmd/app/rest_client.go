package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
)

type RecommRequest struct {
	Tweet struct {
		ID      string   `json:"ID"`
		Text    string   `json:"text"`
		User    string   `json:"user"`
		Request []string `json:"request"`
	} `json:"tweet"`
}

type RecommResponse struct {
	Recommendation struct {
		OriginID  string `json:"originID"`
		Text      string `json:"text"`
		MediaURL  string `json:"mediaURL"`
	} `json:"recommendation"`
	Fulfilled bool
}

func getRecommendation(request RecommRequest) (RecommResponse, error) {
	var recommendation RecommResponse
    var jsonData []byte

    jsonData, err := json.Marshal(request)
	if err != nil {
		return recommendation, fmt.Errorf("Couldn't parse the request to json. Previous Error: %v", err)
    }

    recommAPIResponse, err := postRequest(jsonData)
	if err != nil {
		return recommendation, fmt.Errorf("Couldn't get a recommendation from the remote API. Previous Error: %v", err)
    }
	recommendation.Fulfilled = (recommAPIResponse.StatusCode != 200)
	
	defer recommAPIResponse.Body.Close()

	recommAPIResponseBody, err := ioutil.ReadAll(recommAPIResponse.Body)
    if err != nil {
		return recommendation, fmt.Errorf("Couldn't read API response body. Previous Error: %v", err)
    }

	if err := json.Unmarshal([]byte(recommAPIResponseBody), &recommendation); err != nil {
		return recommendation, fmt.Errorf("Couldn't parse the json response.\nJson: %s \n Previous Error: %v",recommAPIResponseBody, err)
    }

	return recommendation, nil
}

func postRequest(payload []byte) (*http.Response, error){
	client := &http.Client{}

	url := fmt.Sprintf("%s://%s:%s/%s", recommProtocol, recommHost, recommPort, "recommendation")
    
	req, requestErr := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if requestErr != nil {
        return nil, fmt.Errorf("Couldn't create POST request. Previous Error: %v", requestErr)
    }
    req.Header.Add("Content-Type", "application/json")
    
	response, responseErr := client.Do(req)
	if responseErr != nil {
        return nil, fmt.Errorf("Couldn't get a response from recommendation API.\n Got Status Code: %d\n Previous Error: %v",response.StatusCode, requestErr)
    } else if (response.StatusCode != 200 && response.StatusCode != 204){
		return nil, fmt.Errorf("Couldn't get a response from recommendation API.\n Got Status Code: %d\n",response.StatusCode)
	}
	
	return response, nil
}