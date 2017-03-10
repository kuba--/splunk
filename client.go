/*
Package splunk provides the Splunk Enterprise REST API.
More details you can find here: http://dev.splunk.com/restapi
*/
package splunk

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

// A Client is the Splunk REST API client. It keeps credentials and Splunk API endpoints.
// TODO: Add extra endpoints:
// /services/server/control
// /services/server/introspection
// /services/server/logger
// /services/server/roles
// /services/server/settings
// /services/server/status/dispatch-artifacts
// /services/server/status/fishbucket
// /services/server/status/limits
// /services/server/status/partitions-space
type Client struct {
	httpClient *http.Client

	username   string
	password   string
	authUrl    string
	searchUrl  string
	infoUrl    string
	sessionKey string
}

// NewClient returns a new Splunk REST API client.
func NewClient(username string, password string, baseUrl string) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		username:  username,
		password:  password,
		authUrl:   baseUrl + "/services/auth/login",
		searchUrl: baseUrl + "/services/search/jobs/export",
		infoUrl:   baseUrl + "/services/server/info",
	}
}

// Login creates a new session.
func (c *Client) Login(ctx context.Context) error {
	var (
		m  map[string]string
		ok bool
	)

	data := make(url.Values)
	data.Add("username", c.username)
	data.Add("password", c.password)
	data.Add("output_mode", "json")

	resp, err := ctxhttp.PostForm(ctx, c.httpClient, c.authUrl, data)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	if err := json.Unmarshal(body, &m); err != nil {
		return err
	}

	if c.sessionKey, ok = m["sessionKey"]; !ok || c.sessionKey == "" {
		return fmt.Errorf("Login failed: %s\n", string(body))
	}

	return nil
}

// Search streams search results to io.Writer as they become available.
func (c *Client) Search(ctx context.Context, q string, from string, w io.Writer) error {
	data := make(url.Values)
	data.Add("search", fmt.Sprintf("search %s", q))
	data.Add("earliest_time", from)
	data.Add("output_mode", "json")

	req, err := http.NewRequest(http.MethodPost, c.searchUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	if c.sessionKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Splunk %s", c.sessionKey))
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := ctxhttp.Do(ctx, c.httpClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if _, err := w.Write(scanner.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}

// Info streams information to io.Writer about the currently running Splunk instance.
func (c *Client) Info(ctx context.Context, w io.Writer) error {
	data := make(url.Values)
	data.Add("output_mode", "json")

	req, err := http.NewRequest(http.MethodGet, c.infoUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	if c.sessionKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Splunk %s", c.sessionKey))
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := ctxhttp.Do(ctx, c.httpClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if _, err := w.Write(scanner.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}
