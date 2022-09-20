package ftx

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"srkTools"
	"strconv"
	"time"
)

const URL = "https://ftx.com/api/"

type FtxClient struct {
	Client     *http.Client
	Key        string
	Secret     []byte
	Subaccount string
}

func New(key string, secret string, subaccount string) *FtxClient {

	Proxy := db.QueryDB("value", "conf", "name", "PROXY")

	var t *http.Transport
	if Proxy != "" {
		proxy, _ := url.Parse(Proxy)
		t = &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
			Proxy:           http.ProxyURL(proxy),
		}
	} else {
		t = &http.Transport{}
	}
	return &FtxClient{Client: &http.Client{Transport: t}, Key: key, Secret: []byte(secret), Subaccount: url.PathEscape(subaccount)}

}

func Client() *FtxClient {

	key := db.QueryDB("value", "conf", "name", "FTX_KEY")
	secert := db.QueryDB("value", "conf", "name", "FTX_SECERT")

	return New(key, secert, "")
}

func (client *FtxClient) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, client.Secret)
	mac.Write([]byte(signaturePayload))
	//fmt.Println("sign:", hex.EncodeToString(mac.Sum(nil)))
	return hex.EncodeToString(mac.Sum(nil))
}

func (client *FtxClient) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", client.Key)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if client.Subaccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", client.Subaccount)
	}
	//fmt.Println("req:", req)
	return req
}

func (client *FtxClient) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【FTX】res: %s", err))
	}
	return resp, err
}

func (client *FtxClient) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【FTX】res: %s", err))
	}
	return resp, err
}

func (client *FtxClient) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【FTX】res: %s", err))
	}
	return resp, err
}
