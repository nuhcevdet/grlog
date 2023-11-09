package grlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type Level int8

const DebugLevel Level = 0
const InfoLevel Level = 1
const WarnLevel Level = 2
const ErrorLevel Level = 3
const FatalLevel Level = 4
const PanicLevel Level = 5
const NoLevel Level = 6
const Disabled Level = 7
const TraceLevel Level = -1

type gelfMessageStruct struct {
	Version              string                 `json:"version"`
	Host                 string                 `json:"host"`
	ShortMessage         string                 `json:"short_message"`
	FullMessage          string                 `json:"full_message"`
	Timestamp            int64                  `json:"timestamp"`
	LevelID              int                    `json:"level"`
	App                  string                 `json:"_app"`
	Component            string                 `json:"component"`
	Params               map[string]interface{} `json:"params"`
	ipaddress            string                 `json:"-"`
	port                 int                    `json:"-"`
	protocol             string                 `json:"-"`
	errorHandler         func(err error)        `json:"-"`
	connection           net.Conn               `json:"-"`
	alternativeWriteFile bool                   `json:"-"`
}

type Grlog struct {
	graylogIp            string
	grayLogPort          int
	appName              string
	component            string
	protocol             string
	errorHandler         func(err error)
	connection           net.Conn
	alternativeWriteFile bool
	host                 string
}

func (e *Grlog) SetAlternativeLogWriteFile(status bool) {
	e.alternativeWriteFile = status
}

func alternativGrlogFileWriter(message []byte) {
	f, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[Grlog] Log file open error:", err)
	}
	defer f.Close()
	_, err = f.WriteString(string(message))
	if err != nil {
		fmt.Println("[Grlog] File log write error")
	}
	f.WriteString(string("\n"))
}

func (e *Grlog) SetErrorHandler(handler func(err error)) {
	e.errorHandler = handler
}

func (e *Grlog) SetGraylogIp(graylogIp string) {
	e.graylogIp = graylogIp
}

func (e *Grlog) SetGraylogPort(graylogPort int) {
	e.grayLogPort = graylogPort
}

func (e *Grlog) SetProtocol(protocolName string) {
	e.protocol = protocolName
}

func (e *Grlog) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (e *Grlog) SetAppName(appName string) {
	e.appName = appName
}

func (e *Grlog) SetComponentName(componentName string) {
	e.component = componentName
}

func (e *Grlog) SetHostName(hostName string) {
	e.host = hostName
}

func (e *Grlog) New() *gelfMessageStruct {
	return &gelfMessageStruct{
		Version:              "1.1",
		Host:                 e.host,
		Timestamp:            time.Now().Unix(),
		App:                  e.appName,
		Component:            e.component,
		protocol:             e.protocol,
		ipaddress:            e.graylogIp,
		port:                 e.grayLogPort,
		Params:               make(map[string]interface{}),
		errorHandler:         e.errorHandler,
		connection:           e.connection,
		alternativeWriteFile: e.alternativeWriteFile,
	}
}

func (e *gelfMessageStruct) Debug() *gelfMessageStruct {
	e.LevelID = int(DebugLevel)
	return e
}

func (e *gelfMessageStruct) Info() *gelfMessageStruct {
	e.LevelID = int(InfoLevel)
	return e
}

func (e *gelfMessageStruct) Warn() *gelfMessageStruct {
	e.LevelID = int(WarnLevel)
	return e
}

func (e *gelfMessageStruct) Error() *gelfMessageStruct {
	e.LevelID = int(ErrorLevel)
	return e
}

func (e *gelfMessageStruct) Fatal() *gelfMessageStruct {
	e.LevelID = int(FatalLevel)
	return e
}

func (e *gelfMessageStruct) Panic() *gelfMessageStruct {
	e.LevelID = int(PanicLevel)
	return e
}

func (e *gelfMessageStruct) Msg(message string) {
	e.ShortMessage = message
	e.sendLog()
}

func (e *gelfMessageStruct) FullMsg(fullMessage string) *gelfMessageStruct {
	e.FullMessage = fullMessage
	return e
}

func (e *gelfMessageStruct) AddParam(key string, value string) *gelfMessageStruct {
	e.Params[key] = value
	return e
}

func (e *gelfMessageStruct) sendLog() {

	if (e.protocol != "tcp") && (e.protocol != "udp") {
		if e.errorHandler != nil {
			e.errorHandler(errors.New("Protocol is not defined"))
			return
		}
		fmt.Println("Protocol is not defined")
		return
	}
	if e.ipaddress == "" {
		if e.errorHandler != nil {
			e.errorHandler(errors.New("IP address is not defined"))
			return
		}
		fmt.Println("IP address is not defined")
		return
	}
	if e.port == 0 {
		if e.errorHandler != nil {
			e.errorHandler(errors.New("Port is not defined"))
			return
		}
		fmt.Println("Port is not defined")
		return
	}
	logMessage, err := json.Marshal(&e)
	if err != nil {
		if e.errorHandler != nil {
			e.errorHandler(err)
			return
		}
		fmt.Println("Could not marshal log message", err)
		return
	}
	conn, err := net.Dial(e.protocol, fmt.Sprintf("%s:%d", e.ipaddress, e.port))
	if err != nil {
		if e.alternativeWriteFile {
			alternativGrlogFileWriter(logMessage)
			return
		}
		fmt.Println("Could not connect graylog server", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write(logMessage)
	if err != nil {
		if e.alternativeWriteFile {
			alternativGrlogFileWriter(logMessage)
			return
		}
		if e.errorHandler != nil {
			e.errorHandler(err)
			return
		}
		fmt.Println("Could not send log message", err)
		return
	}
}
