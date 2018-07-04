package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/cors"
)

/*
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ike Wai Auth API, %!", r.URL.Path[1:])
}*/

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	const baseURL = "https://agaveauth.its.hawaii.edu"
	const consumerKey = "MYKEY"
	const consumerSecret = "MYSECRET"
	const port = ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"message\":\"Ike Wai Auth\"}"))
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		//get the requests basic auth username and password
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		fmt.Println(pair[0])
		fmt.Println(pair[1])
		//setup request to Agave tenant for fetching an auth token
		//client := &http.Client{}
		data := url.Values{}
		data.Set("grant_type", "password")
		data.Set("username", pair[0])
		data.Set("password", pair[1])
		data.Set("scope", "PRODUCTION")
		data.Set("callbackUrl", "https://localhost:8080")
		fmt.Println(data)
		req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode())) // +"/token?grant_type=password&username="+pair[0]+"&password="+pair[1]+"&callbackUrl=https://localhost:8080", nil) //strings.NewReader(data.Encode()))
		if err != nil {
			//do something
			w.Write([]byte("{\"message\":\"ERROR\"}"))
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(consumerSecret, consumerKey)
		//req.Header.Set("Authorization", "Basic"+basicAuth(consumerKey, consumerSecret))
		fmt.Println(req)
		req.Close = true
		resp, reqErr := http.DefaultClient.Do(req)
		if reqErr != nil {
			fmt.Println(reqErr)
			w.Write([]byte("{\"message\": reqErr}"))
		}

		fmt.Println(resp.Status)
		//fmt.Println(resp.Body)
		//bodyJSON := json.NewDecoder(resp.Body)
		//fmt.Println(bodyJSON)
		//rawbody, _ := ioutil.ReadAll(resp.Body)
		//sfmt.Println(rawbody)
		//jsonBody := json.NewDecoder(resp.Body).Decode(&data)
		//fmt.Println(jsonBody)
		/*var body struct {
			// httpbin.org sends back key/value pairs, no map[string][]string
			Headers map[string]string `json:"headers"`
			Origin  string            `json:"origin"`
		}*/
		processedBody := json.NewDecoder(resp.Body) //.Decode(&body)
		fmt.Println(processedBody)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyJSON := json.NewDecoder(resp.Body)
			fmt.Println(bodyJSON)
			output, _ := json.Marshal(bodyJSON)
			fmt.Println(output)
			w.Header().Set("Content-Type", "application/json")
			w.Write(output)
		}
		//else{
		//	w.Write([]byte(resp.StatusCode))
		//}
		//w.Write([]byte("{\"token\": \"create\"}"))
	})

	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {

		//get the requests basic auth username and password
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)
		fmt.Println(pair[0])
		fmt.Println(pair[1])

		data := url.Values{}
		data.Set("grant_type", "password")
		data.Set("username", pair[0])
		data.Set("password", pair[1])
		data.Set("scope", "PRODUCTION")
		data.Set("callbackUrl", "https://localhost:8080")

		req, err := http.NewRequest("POST", baseURL+"/token", strings.NewReader(data.Encode()))
		if err != nil {
			// handle err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//req.Header.Set("Authorization", "Bearer b7d03a6947b217efb6f3ec3bd3504582")
		req.SetBasicAuth(consumerKey, consumerSecret)
		fmt.Println(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		fmt.Println(resp.Status)
		fmt.Println(resp.Body)
		bodyJSON := json.NewDecoder(resp.Body)
		fmt.Println(bodyJSON)
		type AuthToken struct {
			token string
		}
		NewToken := new(AuthToken)
		jsonResp := json.NewDecoder(resp.Body).Decode(NewToken)
		fmt.Println(jsonResp)
		w.Header().Set("Content-Type", "application/json")
		//w.Write(resp.Body)
		//res, _ := io.Copy(w, resp.Body)
		rbody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rbody))
	})
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	handler := cors.Default().Handler(mux)
	http.ListenAndServeTLS(port, "server.crt", "server.key", handler)

}
