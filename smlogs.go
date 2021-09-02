package SMLogs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Destination      string
	DebugLevel       string
	Flag             string
	DisplayToConsole string
	AppName          string
	AppToken         string
	Setup            bool
}

const (
	Error    = "ERROR"
	Info     = "INFO"
	Critical = "CRITICAL"
	Debug    = "DEBUG"
	Success  = "SUCCESS"
	Ping     = "PING"
)

func (c *Config) New(AppName, AppToken, Destination, DebugLevel, Flag, DisplayToConsole string) error {
	c.AppName = AppName
	c.AppToken = AppToken
	c.Destination = Destination
	c.DebugLevel = DebugLevel
	c.Flag = Flag
	c.DisplayToConsole = DisplayToConsole
	c.Setup = true
	return nil
}

func (c *Config) Error(details ...string) {
	newDetails := append(details, Error)
	go c.Send(newDetails...)
}

func (c *Config) Info(details ...string) {
	newDetails := append(details, Info)
	go c.Send(newDetails...)
}

func (c *Config) Debug(details ...string) {
	newDetails := append(details, Debug)
	go c.Send(newDetails...)
}

func (c *Config) Critical(details ...string) {
	newDetails := append(details, Critical)
	go c.Send(newDetails...)
}

func (c *Config) Success(details ...string) {
	newDetails := append(details, Success)
	go c.Send(newDetails...)
}

func (c *Config) Ping(details ...string) {
	newDetails := append(details, Ping)
	go c.Send(newDetails...)
}

func (c *Config) Send(details ...string) {

	if len(details) < 1 {
		log.Println("Not enough arguments. Pass at least the log contents.")
		return
	}
	//Define vars
	var module string
	var content string
	var status string
	var pkage string

	content = details[0]
	if len(details) >= 2 {
		status = details[1]
	} else {
		status = Error
	}

	if len(details) >= 3 {
		// Get module details
		module = details[2]
		status = details[1]
	} else {
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			lastSlash := strings.LastIndexByte(funcName, '/')
			if lastSlash < 0 {
				lastSlash = 0
			}
			lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash
			pkage = funcName[:lastDot]
			module = funcName[lastDot+1:]

		} else {
			pkage = ""
			module = "NA"
		}

		/*
			_, file, no, ok := runtime.Caller(1)

			if ok {
				module = file + string(no)
			} else  {
				module = "NA"
			}*/
	}

	if !c.Setup {
		log.Println("SMLog is not initialized.")
		return
	}
	if c.DebugLevel != "DEBUG" && status == "DEBUG" {
		return
	}

	//if len(others) > 0 && others[0] > 3 {
	//	log.Println("Returning - total times ", others[0])
	//	return
	//}
	var jsonStr = []byte(`{"contents":"` + jsonEscape(content) + `", "status":"` + jsonEscape(status) + `", "module":"` + jsonEscape(module) + `", "package":"` + jsonEscape(pkage) + `"}`)
	if c.Flag == "Y" {
		if c.DisplayToConsole == "Y" {
			log.Println(bytes.NewBuffer(jsonStr))
		}
		req, err := http.NewRequest("POST", c.Destination, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Debug-Name", c.AppName)
		req.Header.Set("X-Token", c.AppToken)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			/*			totalTimes := 0
						if len(others) > 0 && others[0] > 0 {
							totalTimes = others[0] + 1
						}*/
			time.Sleep(2 * time.Second)
			c.Send("Error in SendInfo() - "+err.Error(), "APPERR", "")
			time.Sleep(1 * time.Second)
			c.Send(content, status, module)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		if resp.StatusCode != 200 {
			log.Println("Debugger Error - response Body:", string(body))
		}
	}
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
