package veritrans

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// get the time after one month from now
func getAfterOneMonth() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(0, 1, 0)
	return expiredAt.Format("01/06")
}

// get the time after one year from now
func getAfterOneYear() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(1, 0, 0)
	return expiredAt.Format("01/06")
}

// process the request
func ProcessRequest(requestURL string, connectionParam *ConnectionParam) (*ConnectionResponse, error) {
	paramByte, err := json.Marshal(connectionParam)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.ParseRequestURI(requestURL)
	if err == nil {
		httpClient := &http.Client{}
		body := bytes.NewBuffer(paramByte)
		req, err := http.NewRequest("POST", parsedURL.String(), body)
		req.Header.Set("Content-Type", "application/json")

		if err == nil {
			res, err := httpClient.Do(req)
			if err == nil {
				body, err := ioutil.ReadAll(res.Body)
				if err == nil {
					var connectionRes ConnectionResponse
					err = json.Unmarshal(body, &connectionRes)
					return &connectionRes, err
				}
			}
		}
	}
	return nil, err
}

// Handler to make hash data of params
func SetHash(connectionParam *ConnectionParam, merchantID, password string) error {
	paramJSON, err := json.Marshal(connectionParam.Params)
	if err != nil {
		return err
	}

	hash := []byte(fmt.Sprintf("%s%s%s", merchantID, paramJSON, password))

	sha := sha256.New()
	sha.Write(hash)
	connectionParam.AuthHash = fmt.Sprintf("%x", sha.Sum(nil))
	return nil
}
