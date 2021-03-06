// Command mkapp creates a mastodon app for a given -instance.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/McKael/madon/v2"
	"within.website/x/internal"
)

var (
	instance    = flag.String("instance", "", "mastodon instance")
	appName     = flag.String("app-name", "Xe/x bot", "app name for mastodon")
	redirectURI = flag.String("redirect-uri", "urn:ietf:wg:oauth:2.0:oob", "redirect URI for app users")
	website     = flag.String("website", "https://within.website/x", "website for users that click the app name")
	username    = flag.String("username", "", "username to generate token")
	password    = flag.String("password", "", "password to generate token")
)

var scopes = []string{"read", "write", "follow"}

func main() {
	internal.HandleStartup()

	c, err := madon.NewApp(*appName, *website, scopes, *redirectURI, *instance)
	if err != nil {
		log.Fatal(err)
	}

	err = c.LoginBasic(*username, *password, scopes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("APP_ID=%s\nAPP_SECRET=%s\nTOKEN=%s\nINSTANCE=%s", c.ID, c.Secret, c.UserToken.AccessToken, *instance)
}
