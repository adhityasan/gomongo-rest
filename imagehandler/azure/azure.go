package azure

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type faceAttr []struct {
	FaceID string `json:"faceId"`
}

// Endpoint of Azure API
type Endpoint struct {
	URL string
	Key string
}

// GetConfidence of two images using Azure Face Verify API
func (e *Endpoint) GetConfidence(jsonParam string) (interface{}, error) {
	res := new(bytes.Buffer)
	if req, err := http.NewRequest("POST", e.URL+"/face/v1.0/verify", strings.NewReader(jsonParam)); err == nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Ocp-Apim-Subscription-Key", e.Key)
		client := &http.Client{}
		resp, errResp := client.Do(req)
		if errResp != nil {
			return nil, errResp
		}
		defer resp.Body.Close()
		res.ReadFrom(resp.Body)
		return res, nil
	}

	return nil, errors.New("Failed to reach Azure Services, may caused by client policy or network connectivity problem")
}

// FaceID get face id from Azure Face Detection API
func (e *Endpoint) FaceID(source interface{}, ch chan interface{}) {
	var res faceAttr
	switch source.(type) {
	case string:
		if req, err := http.NewRequest("POST", e.URL+"/face/v1.0/detect?returnFaceId=true", strings.NewReader(source.(string))); err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Ocp-Apim-Subscription-Key", e.Key)
			client := &http.Client{}
			if resp, err := client.Do(req); err == nil {
				defer resp.Body.Close()
				if err = json.NewDecoder(resp.Body).Decode(&res); err == nil {
					ch <- res[0].FaceID
				}
			}
		}
	default:
		if req, err := http.NewRequest("POST", e.URL+"/face/v1.0/detect?returnFaceId=true", bytes.NewBuffer(source.([]byte))); err == nil {
			req.Header.Set("Content-Type", "application/octet-stream")
			req.Header.Set("Ocp-Apim-Subscription-Key", e.Key)
			client := &http.Client{}
			if resp, err := client.Do(req); err == nil {
				defer resp.Body.Close()
				if err = json.NewDecoder(resp.Body).Decode(&res); err == nil {
					ch <- res[0].FaceID
				}
			}
		}
	}
}
