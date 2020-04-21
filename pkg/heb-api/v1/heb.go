package heb

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	HEB_URL           = "https://www.heb.com"
	LOCATOR_ENDPOINT  = "commerce-api/v1/store/locator/address"
	TIMESLOT_ENDPOINT = "commerce-api/v1/timeslot/timeslots"
)

type Client struct {
	UserAgent string

	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(opts ...ClientOption) *Client {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	c := &Client{httpClient: httpClient}
	c.baseURL, _ = url.Parse(HEB_URL)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type ClientOption func(*Client)

func WithHttpClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBaseURL(u string) func(*Client) {
	return func(c *Client) {
		c.baseURL, _ = url.Parse(u)
	}
}

func (c *Client) LocateStores(zip string, distance int) ([]Store, error) {
	req, err := c.newRequest("POST", LOCATOR_ENDPOINT, map[string]interface{}{"address": zip, "curbsideOnly": false, "radius": distance})
	if err != nil {
		return nil, err
	}

	var lr LocatorResponse

	_, err = c.do(req, &lr)
	if err != nil {
		return nil, err
	}

	var stores []Store
	for _, s := range lr.Stores {
		stores = append(stores, s.Store)
	}

	return stores, nil
}

func (c *Client) GetStoreTimeslots(id string) ([]Timeslot, error) {
	req, err := c.newRequest("GET", TIMESLOT_ENDPOINT, nil, "store_id", id)
	if err != nil {
		return nil, err
	}

	var tr TimeslotResponse

	_, err = c.do(req, &tr)
	if err != nil {
		return nil, err
	}

	var timeslots []Timeslot
	for _, t := range tr.Items {
		timeslots = append(timeslots, t.Timeslot)
	}

	return timeslots, nil

}

func (c *Client) newRequest(method, path string, body interface{}, params ...string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if params != nil {
		q := req.URL.Query()
		for i := 0; i < len(params)-1; i += 2 {
			q.Add(params[i], params[i+1])
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
