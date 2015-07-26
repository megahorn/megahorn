package driver

import (
	"errors"
	"github.com/pusher/pusher-http-go"
	"sync"
	"time"
)

type PusherDriver struct {
	client    *pusher.Client
	channel   string
	event     string
	wg        sync.WaitGroup
	sendQueue chan string
}

func (p *PusherDriver) Configure(config map[string]string) (err error) {
	url, urlOk := config["url"]

	appID, ok := config["app_id"]
	if !ok && !urlOk {
		return errors.New("app_id is required")
	}

	key, ok := config["key"]
	if !ok && !urlOk {
		return errors.New("key is required")
	}

	secret, ok := config["secret"]
	if !ok && !urlOk {
		return errors.New("secret is required")
	}

	secureOpt := true
	secure, ok := config["secure"]
	if ok && secure == "off" {
		secureOpt = false
	}

	host, ok := config["host"]
	if !ok {
		host = ""
	}

	channel, ok := config["channel"]
	if !ok {
		return errors.New("channel is required")
	}

	event, ok := config["event"]
	if !ok {
		return errors.New("event is required")
	}

	if urlOk {
		p.client, err = pusher.ClientFromURL(url)

		if err != nil {
			return
		}
	} else {
		p.client = &pusher.Client{
			AppId:  appID,
			Key:    key,
			Secret: secret,
			Secure: secureOpt,
			Host:   host,
		}
	}

	p.channel = channel
	p.event = event
	p.wg = sync.WaitGroup{}
	p.sendQueue = make(chan string, 32)

	go func() {
		for {
			data := <-p.sendQueue
			p.client.Trigger(p.channel, p.event, map[string]string{"output": data})
			p.wg.Done()
		}
	}()

	return nil
}

func (p *PusherDriver) Write(data []byte) (int, error) {
	p.wg.Add(1)
	p.sendQueue <- string(data)

	return len(data), nil
}

func (p *PusherDriver) Close() error {
	p.wg.Wait()
	return nil
}
