// client.go
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	APIVERSION = "v2"
)

var (
	defaultTimeout = 30 * time.Second
	ErrNotFound    = errors.New("Not found")
)

//Client is ...
type Client struct {
	URL *url.URL

	HTTPClient *http.Client
	// TLSConfig specifies the TLS configuration to use with
	// tls.Client. If nil, the default configuration is used.
	TLSConfig *tls.Config
}

//NewClient return a client.
//You must get a client so that you can access web service.
func NewClient(rawUrl string, tlsConfig *tls.Config) (*Client, error) {
	return NewClientTimeout(rawUrl, tlsConfig, time.Duration(defaultTimeout))
}

func NewClientTimeout(rawUrl string, tlsConfig *tls.Config, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Scheme == "tcp" {
		if tlsConfig == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}
	httpClient := newHTTPClient(u, tlsConfig, timeout)
	return &Client{u, httpClient, tlsConfig}, nil
}

func newHTTPClient(u *url.URL, tlsConfig *tls.Config, timeout time.Duration) *http.Client {
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	switch u.Scheme {
	case "unix":
		socketPath := u.Path
		unixDial := func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, timeout)
		}
		httpTransport.Dial = unixDial
		//Override the main URL object so the HTTP lib won't complain
		u.Scheme = "http"
		u.Host = "unix.sock"
		u.Path = ""
	default:
		httpTransport.Dial = func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout(proto, addr, timeout)
		}
	}
	return &http.Client{Transport: httpTransport}
}

func (client *Client) doRequest(method string, path string, body []byte, headers map[string]string) ([]byte, error) {
	b := bytes.NewBuffer(body)

	reader, err := client.doStreamRequest(method, path, b, headers)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (client *Client) doStreamRequest(method string, path string, in io.Reader, headers map[string]string) (io.ReadCloser, error) {
	if (method == "POST" || method == "PUT") && in == nil {
		in = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, client.URL.String()+path, in)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "connection refused") && client.TLSConfig == nil {
			return nil, fmt.Errorf("%v. Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

func (client *Client) Catelog() (*Repositories, error) {
	uri := fmt.Sprintf("/%s/_catalog", APIVERSION)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}
	var repositories Repositories
	if err = json.Unmarshal(data, &repositories); err != nil {
		return nil, err
	}
	return &repositories, nil
}

func (client *Client) Tags(name string) (*Tags, error) {
	uri := fmt.Sprintf("/%s/%s/tags/list", APIVERSION, name)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}
	var tags Tags
	err = json.Unmarshal(data, &tags)
	if err != nil {
		return nil, err
	}
	return &tags, nil
}
