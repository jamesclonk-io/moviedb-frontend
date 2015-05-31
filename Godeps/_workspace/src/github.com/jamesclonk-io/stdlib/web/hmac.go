package web

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (bc *BackendClient) addHmacHeader(method, backendURL, data string, req *http.Request) error {
	u, err := url.Parse(backendURL)
	if err != nil {
		return err
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	payload, err := hmacPayload(u, method, data, timestamp)
	if err != nil {
		return err
	}
	signature, err := hmacSignature(payload, bc.secret)
	if err != nil {
		return err
	}

	req.Header.Add("X-Jcio-Timestamp", timestamp)
	req.Header.Add("X-Jcio-Hmac", base64.StdEncoding.EncodeToString(signature))

	return nil
}

func (b *Backend) hmacAuth(req *http.Request) (bool, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, err
	}

	timestamp := req.Header.Get("X-Jcio-Timestamp")
	if len(timestamp) == 0 {
		return false, nil
	}

	// check how much time has passed since 'timestamp' to prevent replay attacks
	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, err
	}
	if time.Now().Sub(time.Unix(unix, 0)).Seconds() > 10 {
		return false, nil
	}

	// calculate expected hmac signature
	payload, err := hmacPayload(req.URL, req.Method, string(data), timestamp)
	if err != nil {
		return false, err
	}
	expectedHmacSignature, err := hmacSignature(payload, b.secret)
	if err != nil {
		return false, err
	}
	givenHmacSignature, err := base64.StdEncoding.DecodeString(req.Header.Get("X-Jcio-Hmac"))
	if err != nil {
		return false, err
	}

	// compare expected and given hmac signatures
	return hmac.Equal(givenHmacSignature, expectedHmacSignature), nil
}

func hmacPayload(url *url.URL, method, data, timestamp string) ([]byte, error) {
	payload := make(map[string]string)
	payload["method"] = method
	payload["uri"] = url.RequestURI()
	payload["query"] = url.RawQuery
	payload["data"] = data
	payload["timestamp"] = timestamp

	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func hmacSignature(payload, secret []byte) ([]byte, error) {
	h := hmac.New(sha512.New, secret)
	_, err := h.Write(payload)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
