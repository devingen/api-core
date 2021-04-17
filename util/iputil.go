package util

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// GetClientIPHelper gets the client IP using a mixture of techniques.
// This is how it is with golang at the moment.
func GetClientIPHelper(req *http.Request) (ipResult string, errResult error) {

	// Try lots of ways :) Order is important.

	//  Try Request Header ("Origin")
	url, err := url.Parse(req.Header.Get("Origin"))
	if err == nil {
		host := url.Host
		ip, _, err := net.SplitHostPort(host)
		if err == nil {
			return ip, nil
		}
	}

	// Try by Request
	ip, err := getClientIPByRequestRemoteAddr(req)
	if err == nil {
		return ip, nil
	}

	// Try Request Headers (X-Forwarder). Client could be behind a Proxy
	ip, err = getClientIPByHeaders(req)
	if err == nil {
		return ip, nil
	}

	err = errors.New("error: Could not find clients IP address")
	return "", err
}

// getClientIPByRequest tries to get directly from the Request.
// https://blog.golang.org/context/userip/userip.go
func getClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {

	// Try via request
	ip, _, err = net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
		return "", fmt.Errorf(message)

	}
	return userIP.String(), nil

}

// getClientIPByHeaders tries to get directly from the Request Headers.
// This is only way when the client is behind a Proxy.
func getClientIPByHeaders(req *http.Request) (ip string, err error) {

	// Client could be behid a Proxy, so Try Request Headers (X-Forwarder)
	ipSlice := []string{}

	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		if v != "" {
			return v, nil
		}
	}
	err = errors.New("error: Could not find clients IP address from the Request Headers")
	return "", err

}
