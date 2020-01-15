package logger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fluent/fluent-logger-golang/fluent"

	"github.com/golang-common-packages/template/model"
)

// FluentClient manage all fluent action
type FluentClient struct {
	Config *model.Fluent
	Client *fluent.Fluent
}

/*
	@sessionMapping: Mapping between model.Fluent and FluentClient for singleton pattern
*/
var (
	sessionMapping = make(map[string]*FluentClient)
)

// NewFluentClient function return a new fluent client based on singleton pattern
func NewFluentClient(config *model.Service) LoggerStore {
	hash := config.Hash()
	currentSession := sessionMapping[hash]
	if currentSession == nil {
		currentSession = &FluentClient{nil, nil}

		logger, err := fluent.New(getConfig(config.Fluent))
		if err != nil {
			log.Println("Error when try to connect to Fluent server: ", err)
			panic(err)
		}
		log.Println("Connected to Fluent Server")

		currentSession = &FluentClient{&config.Fluent, logger}
		sessionMapping[hash] = currentSession
	}
	return currentSession
}

// Write
func (c *FluentClient) Write(p []byte) (n int, err error) {
	data := make(map[string]interface{})
	err = json.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}

	if c.Client == nil {
		fmt.Print(string(p))
		return 0, err
	}

	err = c.Client.Post(c.Config.Tag, data)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

// Close
func (c *FluentClient) Close() {
	c.Client.Close()
}

// getConfig function return config of fluent
func getConfig(f model.Fluent) fluent.Config {
	return fluent.Config{
		FluentPort:         f.Port,
		FluentHost:         f.Host,
		TagPrefix:          f.Prefix,
		MarshalAsJSON:      false,
		SubSecondPrecision: true,
	}
}
