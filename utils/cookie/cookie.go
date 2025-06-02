// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cookie provides cookie related functions
package cookie

import (
	"errors"
	"net/http"
)

// GetCookie  get cookie by name
func GetCookie(rawCookies, name string) (*http.Cookie, error) {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	cookie, err := request.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}
	return cookie, nil
}

// Parse  Cookie to []*http.Cookie
func Parse(rawCookies string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}
	return request.Cookies()
}

// Format  Cookie to string
func Format(cookies []*http.Cookie) []string {
	rawCookies := make([]string, len(cookies))
	for _, item := range cookies {
		rawCookies = append(rawCookies, item.String())
	}
	return rawCookies
}
