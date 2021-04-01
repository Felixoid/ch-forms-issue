// Package main provides issue reproducing
package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func buildBody(lines int, u *url.URL) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	header := ""
	writer := multipart.NewWriter(body)
	name := "metrics_list"
	part, err := writer.CreateFormFile(name, name)
	if err != nil {
		return nil, header, err
	}
	data := []byte{}

	for i := 1; i <= lines; i++ {
		data = append(data, []byte("some pretty long line to fill the post\n")...)
	}

	// Send each table in separated form
	_, err = part.Write(data)
	if err != nil {
		return nil, header, err
	}

	// Set name_format and name_structure for the table
	q := u.Query()
	q.Set(name+"_format", "TSV")
	q.Set(name+"_structure", "Path String")
	u.RawQuery = q.Encode()

	err = writer.Close()
	if err != nil {
		return nil, header, err
	}

	header = writer.FormDataContentType()
	return body, header, nil
}

func main() {
	dsn := flag.String("url", "https://localhost:8443", "url to make request")
	lines := flag.Int("lines", 40, "number of lines 'some pretty long line to fill the post' in request")
	withAgent := flag.Bool("agent", false, "if set, custom agent will be used")
	flag.Parse()
	url, err := url.Parse(*dsn)
	if err != nil {
		return
	}
	q := url.Query()
	q.Set("query", `SELECT
          Path,
          arrayFilter(x->isNotNull(x), anyOrNullResample(1615852800, 1617235199, 86400)(toUInt32(intDiv(Time, 86400)*86400), Time)),
          arrayFilter(x->isNotNull(x), avgOrNullResample(1615852800, 1617235199, 86400)(Value, Time))
        FROM default.test
        PREWHERE Date >='2021-03-16' AND Date <= '2021-04-01'
        WHERE (Path in metrics_list) AND (Time >= 1615852800 AND Time <= 1617235199)
        GROUP BY Path FORMAT TSV`)
	url.RawQuery = q.Encode()
	postBody, contentHeader, err := buildBody(*lines, url)
	req, err := http.NewRequest("POST", url.String(), postBody)
	if *withAgent {
		req.Header.Add("User-Agent", "Graphite-Clickhouse/0.12.0 (table:graphite.data)")
	}
	req.Header.Add("Content-Type", contentHeader)
	req.Header.Add("Content-Encoding", "gzip")
	client := &http.Client{
		Timeout: time.Second * 50,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * 3,
			}).Dial,
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		},
	}
	rawRequest, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		fmt.Print(string(rawRequest))
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Print(resp)
}
