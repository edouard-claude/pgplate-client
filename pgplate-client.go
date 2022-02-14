package pgplateclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	BaseUrl      string
	OAuthID      string
	OAuthSecret  string
	accessToken  *string
	Jwt          *string
	RefreshToken *string
	client       *http.Client
}

type methodHttp string

const (
	GET    methodHttp = "GET"
	POST   methodHttp = "POST"
	PUT    methodHttp = "PUT"
	PATCH  methodHttp = "PATCH"
	DELETE methodHttp = "DELETE"
)

func NewClient(client *Client) (*Client, error) {

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: customTransport,
	}

	if err := client.getAuthorization(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Fetch(route string, method methodHttp, payload []byte, response interface{}) ([]byte, error) {

	req, err := http.NewRequest(string(method), c.BaseUrl+route, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	if c.accessToken != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *c.accessToken))
	} else {
		c.getAuthorization()
	}

	if c.Jwt != nil {
		req.Header.Add("jwtToken", fmt.Sprintf("Bearer %s", *c.Jwt))
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		c.getAuthorization()
		c.Fetch(route, method, payload, response)
	}

	json.NewDecoder(res.Body).Decode(response)

	return nil, nil
}

func (c *Client) getAuthorization() error {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	writer.WriteField("grant_type", "client_credentials")
	writer.WriteField("client_id", c.OAuthID)
	writer.WriteField("client_secret", c.OAuthSecret)
	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseUrl+"/authorization", payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	target := &response{}
	json.NewDecoder(res.Body).Decode(target)

	c.accessToken = &target.AccessToken

	return nil
}
