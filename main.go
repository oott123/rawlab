package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		parts := strings.SplitN(request.URL.Path, "@", 2)
		repoName := url.QueryEscape(parts[0][1:])
		fParts := strings.SplitN(parts[1], "/", 2)
		refName := url.QueryEscape(fParts[0])
		fileName := url.QueryEscape(fParts[1])
		fileExt := filepath.Ext(fParts[1])
		fileMime := mime.TypeByExtension(fileExt)
		if fileExt == ".ts" {
			fileMime = "application/typescript"
		}

		authToken := request.URL.Query().Get("token")
		if authToken == "" {
			authHeader := request.Header.Get("authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				authToken = parts[len(parts)-1]
			}
		}

		baseUrl := os.Getenv("GITLAB_API")
		requestUrl := fmt.Sprintf("%s/projects/%s/repository/files/%s/raw?ref=%s", baseUrl, repoName, fileName, refName)

		apiReq, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write(([]byte)(err.Error()))
			return
		}

		apiReq.Header.Add("User-Agent", "RawLab/1.0")
		if authToken != "" {
			apiReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
		}
		request.Header.Del("User-Agent")
		request.Header.Del("Authorization")
		request.Header.Del("X-Forwarded-For")
		request.Header.Del("X-Real-IP")
		request.Header.Del("CF-Connecting-IP")
		request.Header.Del("Keep-Alive")
		request.Header.Del("Transfer-Encoding")
		request.Header.Del("TE")
		request.Header.Del("Connection")
		request.Header.Del("Trailer")
		request.Header.Del("Upgrade")
		request.Header.Del("Proxy-Authorization")
		request.Header.Del("Proxy-Authenticate")
		request.Header.Del("Cookie")
		for k, l := range request.Header {
			apiReq.Header.Del(k)
			for _, v := range l {
				apiReq.Header.Add(k, v)
			}
		}

		apiRes, err := http.DefaultClient.Do(apiReq)
		if err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write(([]byte)(err.Error()))
			return
		}

		apiReq.Header.Del("Set-Cookie")
		apiReq.Header.Del("Keep-Alive")
		apiReq.Header.Del("Transfer-Encoding")
		apiReq.Header.Del("TE")
		apiReq.Header.Del("Connection")
		apiReq.Header.Del("Trailer")
		apiReq.Header.Del("Upgrade")
		if apiRes.StatusCode == 200 {
			apiRes.Header.Del("Content-Type")
			apiRes.Header.Add("Content-Type", fileMime)
		}

		for k, v := range apiRes.Header {
			writer.Header().Del(k)
			for _, h := range v {
				writer.Header().Add(k, h)
			}
		}
		writer.WriteHeader(apiRes.StatusCode)

		io.Copy(writer, apiRes.Body)
	})

	listen := ":8625"

	if os.Getenv("PORT") != "" {
		listen = fmt.Sprintf(":%s", os.Getenv("PORT"))
	} else if os.Getenv("LISTEN") != "" {
		listen = os.Getenv("LISTEN")
	}

	go func() {
		fmt.Printf("Server listened on %s\n", listen)
	}()
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err)
	}
}
