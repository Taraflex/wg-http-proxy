package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
func toJson [T any](i T) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
*/
type cData struct {
	Content string
	Time    time.Time
}

func trimEtag(s string) string {
	return strings.Trim(s, `W/"`) // для случая W/"W...W" коллизии не произойдет так как etag имеет фиксированную длинну
}

func ProxyHandler(client http.Client) func(http.ResponseWriter, *http.Request) {
	var cacheMutex sync.Mutex
	cache := make(map[string]cData)

	contentContainer := [1024 * 1025]byte{} // 1MB + 1KB

	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/health" {
			w.WriteHeader(204)
		} else if req.URL.Path == "/PKH.pac" {
			ctx := req.Context()

			var err error
			etag := ""
			noneMatch := trimEtag(req.Header.Get("If-None-Match"))

			defer func() {
				if ctx.Err() == nil {
					if err != nil {
						log.Println(err)
					}
					cached := cache[etag]
					if noneMatch != "" && cached.Content == "" {
						etag = noneMatch
						cached = cache[etag]
					}
					if !cached.Time.IsZero() && (etag == noneMatch || cached.Content != "") {
						h := w.Header()
						h.Set("Cache-Control", "public, max-age=86400")
						//h.Set("Content-Type", "text/javascript; charset=utf-8")
						h.Set("Content-Type", "application/x-ns-proxy-autoconfig; charset=utf-8")
						h.Set("Date", cached.Time.Format(http.TimeFormat))
						h.Set("Etag", `W/"`+etag+`"`)
						if etag == noneMatch {
							w.WriteHeader(http.StatusNotModified)
						} else {
							w.WriteHeader(http.StatusOK)
							io.WriteString(w, cached.Content)
						}
					} else if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}()

			azRequest, err := http.NewRequestWithContext(ctx, "GET", "http://antizapret.prostovpn.org/proxy.pac", nil)
			if err != nil {
				return
			}
			headers := req.Header.Clone()

			reqEtag := noneMatch
			if reqEtag == "" {
				var lastDate time.Time
				for tag, data := range cache {
					if data.Content != "" && data.Time.After(lastDate) {
						lastDate = data.Time
						reqEtag = tag
					}
				}
			}
			if reqEtag == "" {
				headers.Set("Cache-Control", "public, max-age=0")
			} else {
				headers.Set("Cache-Control", "public, max-age=31536000")
				headers.Set("If-None-Match", `W/"`+reqEtag+`"`)
			}
			// wget -d http://www.google.com -O/dev/null 2>&1 |grep ^User-Agent
			headers.Set("User-Agent", "Wget/1.21")

			azRequest.Header = headers
			azResponse, err := client.Do(azRequest)
			if err != nil {
				return
			}
			defer azResponse.Body.Close()

			etag = trimEtag(azResponse.Header.Get("Etag"))

			if azResponse.StatusCode == http.StatusNotModified || (etag != "" && etag == noneMatch) {
				fmt.Printf("etag: %v\n", etag)
			} else if azResponse.StatusCode != http.StatusOK {
				fmt.Printf("response status: %v\n", azResponse.Status)
				err = errors.New(azResponse.Status)
				return
			}

			if etag != "" {
				func() {
					t, _ := http.ParseTime(azResponse.Header.Get("Date"))
					if t.IsZero() {
						t = time.Now().Add(-time.Minute)
					}

					cacheMutex.Lock()
					defer cacheMutex.Unlock()

					prevContent := cache[etag].Content
					if etag != noneMatch && prevContent == "" {
						var source io.Reader = azResponse.Body
						//n, err = io.ReadFull(azResponse.Body, contentContainer[:])
						n := 0
						for n < 1024*1025 && err == nil {
							var nn int
							nn, err = source.Read(contentContainer[n:])
							n += nn
							log.Print(n)
						}
						if err == io.EOF {
							err = nil
						}

						if err == nil {
							pacData := contentContainer[:n]
							cache[etag] = cData{
								Content: string(pacData), //todo patch proxy uri
								//direct := mapA(strings.Split(req.URL.Query().Get("direct"), ","), strings.TrimSpace)
								//wg := mapA(strings.Split(req.URL.Query().Get("wg"), ","), strings.TrimSpace)
								Time: t,
							}
							return
						}
					}

					if cache[etag].Time.Before(t) {
						cache[etag] = cData{
							Content: prevContent,
							Time:    t,
						}
					}
				}()
			}
		} else {
			http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", http.StatusInternalServerError)
		}
	}
}
