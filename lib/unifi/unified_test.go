package unifi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
	//	log "github.com/Sirupsen/logrus"
)

var (
	ctx = context.TODO()
)

func setup() {
	client := NewUniFiClient(nil, nil)
	value, err := url.Parse("https://192.168.10.7:8443/api/s/")

	if err != nil {
		//log.Fatal(err)
	} else {
		client.BaseURL = value
	}
}

func TestNewRequest(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: time.Second * 300, Transport: tr}
	_, err := client.Get("https://golang.org/")
	if err != nil {
		fmt.Println(err)
	}

	d := &UnifiedDBOptions{
		DbUsageEnabled: true,
		UseInMemoryDB:  true,
	}

	o := &UnifiedOptions{
		DbUsage: d,
	}

	c := NewUniFiClient(client, o)

	c.Authentication.Login(ctx, "donstewa", "Blackadder!782")
	alarms, _, err := c.Alarms.List(ctx, nil)
	fmt.Println(len(alarms), " alarms returned.")
	//log.Info("Alarms Response : ", resp)
	//log.WithFields(log.Fields{
	//	"json": alarms,
	//}).Info("Alarms found.")
	events, _, err := c.Events.List(ctx, nil)
	fmt.Println(len(events), " events returned.")
	//log.Info("Events Response : ", resp)
	//log.WithFields(log.Fields{
	//	"json": events,
	//}).Info("Events found.")
	users, _, err := c.Users.List(ctx, nil)
	fmt.Println(len(users), " users returned.")
	//log.Info("Users Response : ", resp)
	//log.WithFields(log.Fields{
	//	"json": users,
	//}).Info("Users found.")

	c.Stop()
}
