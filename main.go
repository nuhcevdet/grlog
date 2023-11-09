package main

import (
	"fmt"
	"grlog/grlog"
	"time"
)

func errorHandler(err error) {
	fmt.Println("Error", err)
}

func main() {
	grlog := grlog.Grlog{}
	grlog.SetAlternativeLogWriteFile(true)
	grlog.SetErrorHandler(errorHandler)
	grlog.SetGraylogIp("172.22.29.1")
	grlog.SetGraylogPort(12201)
	grlog.SetProtocol("tcp")
	grlog.SetAppName("DenemeApp")
	grlog.SetComponentName("DenemeCompopnent")
	for range time.Tick(1 * time.Second) {
		grlog.New().Info().Msg(fmt.Sprintf("%d Test", time.Now().Second()))
	}
}
