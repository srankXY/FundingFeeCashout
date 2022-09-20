package gate

import (
	"FundingFeeCashout/db"
	"FundingFeeCashout/lib"
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"srkTools"
	"strconv"
	"time"
)

type GateClient struct {
	Client *http.Client
	Key    string
	Secret string
}

const URL = "https://api.gateio.ws/api/v4"

func New(key string, secret string) *GateClient {

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
	return &GateClient{Client: &http.Client{Transport: t}, Key: key, Secret: secret}
}

func Client() *GateClient {

	key := db.QueryDB("value", "conf", "name", "GATE_KEY")
	secert := db.QueryDB("value", "conf", "name", "GATE_SECERT")

	return New(key, secert)
}

func (client *GateClient) sign(signaturePayload string) string {
	mac := hmac.New(sha512.New, []byte(client.Secret))
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (client *GateClient) hashBody(body []byte) string {
	hBody := sha512.New()
	hBody.Write(body)
	return hex.EncodeToString(hBody.Sum(nil))
}

func (client *GateClient) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	requestUrl, _ := url.Parse(path)
	//query := requestUrl.Query()
	//requestUrl.RawQuery = query.Encode()
	//rawQuery, _ := url.QueryUnescape(requestUrl.RawQuery)

	shaBody := client.hashBody(body)
	signaturePayload := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, "/api/v4"+requestUrl.Path, requestUrl.RawQuery, shaBody, ts)
	//signaturePayload := "GET\n/api/v4/futures/orders\ncontract=BTC_USD&status=finished&limit=50\ncf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e\n1541993715"
	signature := client.sign(signaturePayload)

	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", client.Key)
	req.Header.Set("SIGN", signature)
	req.Header.Set("Timestamp", ts)
	//fmt.Println("【GATE】req:", req)
	return req
}

func (client *GateClient) _get(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("GET", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【GATE】res: %s", err))
	}
	return resp, err
}

func (client *GateClient) _post(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("POST", path, body)
	resp, err := client.Client.Do(preparedRequest)
	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【GATE】res: %s", err))
	}
	return resp, err
}

func (client *GateClient) _delete(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("DELETE", path, body)
	resp, err := client.Client.Do(preparedRequest)

	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【GATE】res: %s", err))
	}

	return resp, err
}

func (client *GateClient) _put(path string, body []byte) (*http.Response, error) {
	preparedRequest := client.signRequest("PUT", path, body)
	resp, err := client.Client.Do(preparedRequest)

	if err != nil {
		srkTools.DebugLog(lib.DebugLevel.VERBOSE, fmt.Sprintf("【GATE】res: %s", err))
	}

	return resp, err
}
