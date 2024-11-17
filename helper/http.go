package helper

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"time"
)

type PHttp struct {
	Timeout            time.Duration
	KeepAlive          time.Duration
	IsDisableKeepAlive bool
	MaxIdleConns       int
	IdleConnTimeout    time.Duration
	DisableCompression bool
}

// HttpDial (string)
func HttpDial(url string, t time.Duration) error {
	timeout := t * time.Second
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		fmt.Printf("Site unreachable : %s, error: %#v\n", url, err)
	} else {
		defer conn.Close()
	}

	return err
}

func HttpDial2(url string, t time.Duration) bool {
	_, err := http.Get(url)
	if err != nil {
		return false
		//fmt.Printf("%#v", err.Error())
	} else {
		return true
		//fmt.Printf("%d, %s", resp.StatusCode, resp.Status)
	}
}

// HttpClient (time.Duration, time.Duration, bool)
func HttpClient(p PHttp) *http.Client {
	//ref: Copy and modify defaults from https://golang.org/src/net/http/transport.go
	//Note: Clients and Transports should only be created once and reused
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			// Modify the time to wait for a connection to establish
			Timeout:   1 * time.Second,
			KeepAlive: p.KeepAlive * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   p.IsDisableKeepAlive,
		MaxIdleConns:        p.MaxIdleConns,
		IdleConnTimeout:     p.IdleConnTimeout,
		DisableCompression:  p.DisableCompression,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   p.Timeout * time.Second,
	}

	return &client
}

// @return param : respBody, response.Status, response.StatusCode, strlog, err
func Get(url string, headers map[string]string, transport PHttp) ([]byte, string, int, string, error) {

	start := time.Now()

	var (
		respBody    []byte
		elapseInSec string
		elapseInMS  string
		strlog      string
	)

	httpClient := HttpClient(transport)

	req, err := http.NewRequest("GET", url, nil)
	if len(headers) != 0 {
		for k, v := range headers {

			if k == "Basic-Auth" {
				auth := strings.Split(v, ":")
				req.SetBasicAuth(auth[0], auth[1])
			} else {
				req.Header.Set(k, v)
			}
		}
	}
	req.Close = true

	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Error Occured : %#v", err), err
	}

	var (
		getConn   string
		dnsStart  string
		dnsDone   string
		connStart string
		connDone  string
		gotConn   string
	)

	clientTrace := &httptrace.ClientTrace{
		GetConn:  func(hostPort string) { getConn = fmt.Sprintf("Starting to create conn : [%s], ", hostPort) },
		DNSStart: func(info httptrace.DNSStartInfo) { dnsStart = fmt.Sprintf("starting to look up dns : [%s], ", info) },
		DNSDone:  func(info httptrace.DNSDoneInfo) { dnsDone = fmt.Sprintf("done looking up dns : [%#v], ", info) },
		ConnectStart: func(network, addr string) {
			connStart = fmt.Sprintf("starting tcp connection : [%s, %s], ", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			connDone = fmt.Sprintf("tcp connection created [%s, %s, %#v], ", network, addr, err)
		},
		GotConn: func(info httptrace.GotConnInfo) { gotConn = fmt.Sprintf("conn was reused: [%#v]", info) },
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	response, err := httpClient.Do(req)
	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Error sending request to API endpoint : %#v, Hit: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", err, url, string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)), err
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	respBody, err = io.ReadAll(response.Body)
	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Couldn't parse response body : %#v, Hit: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", err, url, string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)), err
	}

	elapse := time.Since(start)

	elapseInSec = fmt.Sprintf("%f", elapse.Seconds())
	elapseInMS = strconv.FormatInt(elapse.Milliseconds(), 10)

	strlog = fmt.Sprintf("Hit: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", url, string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn))

	req = nil
	httpClient = nil

	return respBody, response.Status, response.StatusCode, strlog, err
}

// @return param : respBody, response.Status, response.StatusCode, strlog, err
func Post(url string, headers map[string]string, body []byte, transport PHttp) ([]byte, string, int, string, error) {

	start := time.Now()

	var (
		respBody    []byte
		elapseInSec string
		elapseInMS  string
		strlog      string
	)

	httpClient := HttpClient(transport)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	//req.Header.Set("Content-Type", content_type)

	if len(headers) != 0 {
		for k, v := range headers {

			if k == "Basic-Auth" {
				auth := strings.Split(v, ":")
				req.SetBasicAuth(auth[0], auth[1])
			} else {
				req.Header.Set(k, v)
			}
		}
	}

	req.Close = true

	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Error Occured : %#v", err), err
	}

	var (
		getConn   string
		dnsStart  string
		dnsDone   string
		connStart string
		connDone  string
		gotConn   string
	)

	clientTrace := &httptrace.ClientTrace{
		GetConn:  func(hostPort string) { getConn = fmt.Sprintf("Starting to create conn : [%s], ", hostPort) },
		DNSStart: func(info httptrace.DNSStartInfo) { dnsStart = fmt.Sprintf("starting to look up dns : [%s], ", info) },
		DNSDone:  func(info httptrace.DNSDoneInfo) { dnsDone = fmt.Sprintf("done looking up dns : [%#v], ", info) },
		ConnectStart: func(network, addr string) {
			connStart = fmt.Sprintf("starting tcp connection : [%s, %s], ", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			connDone = fmt.Sprintf("tcp connection created [%s, %s, %#v], ", network, addr, err)
		},
		GotConn: func(info httptrace.GotConnInfo) { gotConn = fmt.Sprintf("conn was reused: [%#v]", info) },
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)

	response, err := httpClient.Do(req)
	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Error sending request to API endpoint : %#v, Hit: %s, Request: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", err, url, string(body), string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)), err
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	respBody, err = io.ReadAll(response.Body)
	if err != nil {

		return []byte(""), "", 0, fmt.Sprintf("Couldn't parse response body : %#v, Hit: %s, Request: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", err, url, string(body), string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn)), err
	}

	elapse := time.Since(start)

	elapseInSec = fmt.Sprintf("%f", elapse.Seconds())
	elapseInMS = strconv.FormatInt(elapse.Milliseconds(), 10)

	strlog = fmt.Sprintf("Hit: %s, Request: %s, Response: %s, Status: %s, Status Code: %d, Elapse: %s second, %s milisecond, live trace : %s", url, string(body), string(respBody), response.Status, response.StatusCode, elapseInSec, elapseInMS, Concat(getConn, dnsStart, dnsDone, connStart, connDone, gotConn))

	req = nil
	httpClient = nil

	return respBody, response.Status, response.StatusCode, strlog, err
}
