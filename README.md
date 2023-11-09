
# GRLOG - Golang Graylog Log Package

This project was created to send logs to the graylog application in a simple way.

With this project, you can create logs in accordance with Graylog. Log levels are compatible with graylog, and applications and components can be entered to classify logs. When a connection to Graylog cannot be established or the log cannot be sent, the logs are set to be written to the file.

# Log Levels
Debug Level

Info Level

Warn Level

Error Level

Fatal Level

Panic Level

No Level

Disabled

Trace Level


# Example

```
grlog := grlog.Grlog{}
//Should the file be written to when the connection to graylog is lost?
grlog.SetAlternativeLogWriteFile(true)
//An error handler function to catch errors encountered
grlog.SetErrorHandler(errorHandler)
//Graylog ip address
grlog.SetGraylogIp("127.0.0.1")
//Graylog port number
grlog.SetGraylogPort(12201)
//graylog log send protocol udp or tcp
grlog.SetProtocol("tcp")
//Graylog custom field App
grlog.SetAppName("TestApp")
//Graylog custom field Component
grlog.SetComponentName("TestComponent")
for range time.Tick(1 * time.Second) {
	grlog.New().Info().Msg(fmt.Sprintf("%d Test", time.Now().Second()))
}
```
