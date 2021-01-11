package utils

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
)

//ErrorReport is a struct for marshalling errors to send to discord.
type ErrorReport struct {
	Err string `json:"content"`
}

//Errors is a rolling counter for the number of errors that have been encountered, but weren't posted to discord because of connection/other issues
var Errors int

//LogError sends an error to discord. If it fails, it logs it locally to a json log file.
func LogError(logErr string) {
	errStruct := ErrorReport{Err: logErr}

	b, err := json.Marshal(errStruct)
	if err != nil {
		logLocalError(err.Error())
		logLocalError(logErr)

		return
	}

	_, err = net.LookupIP("google.com") //The domain used here could be anything, as it's only used for testing DNS
	if err != nil {
		logLocalError("DNS lookup failed. Logging error locally...")
		logLocalError(logErr)

		return
	}

	resp, err := http.Post(Config.DiscordWebhook, "application/json", bytes.NewReader(b)) //Discord webhook post
	if err != nil {
		logLocalError(err.Error())
		logLocalError(logErr)

		return
	}
	defer resp.Body.Close()

	return
}

func logLocalError(logErr string) {
	Errors++

	f, err := os.OpenFile("./error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(logErr))
	if err != nil {
		panic(err)
	}
}
