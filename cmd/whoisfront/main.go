// Command whoisfront is a simple CGI wrapper to switchcounter.science. This is used in some internal tooling.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"

	"within.website/x/internal"
	"within.website/x/web/switchcounter"
)

var (
	switchCounterURL = flag.String("switch-counter-url", "", "the webhook for switchcounter.science")

	sc switchcounter.API
)

func main() {
	internal.HandleStartup()

	sc = switchcounter.NewHTTPClient(*switchCounterURL)

	err := cgi.Serve(http.HandlerFunc(handle))
	if err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		front, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		req := sc.Switch(string(front))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		err = switchcounter.Validate(resp)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, string(front))
		return
	}

	req := sc.Status()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	err = switchcounter.Validate(resp)
	if err != nil {
		panic(err)
	}
	var st switchcounter.Status
	err = json.NewDecoder(resp.Body).Decode(&st)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, st.Front)
}
