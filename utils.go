package veritrans

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// get the time after one month from now
func GetAfterOneMonth() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(0, 1, 0)
	return expiredAt.Format("01/06")
}

// get the time after one year from now
func GetAfterOneYear() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(1, 0, 0)
	return expiredAt.Format("01/06")
}

// get random id
func GetRandomID(digit int) int {
	rand.Seed(time.Now().UnixNano())
	low := int(math.Pow10(digit - 1))
	high := int(math.Pow10(digit) - 1)
	return low + rand.Intn(high-low)
}

// process the request
func ProcessRequest(requestURL string, connectionParam *ConnectionParam) (*ConnectionResponse, error) {
	var err error
	paramByte, err := json.Marshal(connectionParam)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	body := bytes.NewBuffer(paramByte)
	req, err := http.NewRequest("POST", parsedURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var connectionRes ConnectionResponse
	err = json.Unmarshal(resBody, &connectionRes)
	return &connectionRes, err
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
