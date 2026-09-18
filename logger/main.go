package main

import "fmt"

func main() {
	fmt.Println("hi we are designing a logger system")

	fileWriter, _ := NewFileWriter("app.log")

	log := NewLoggerBuilder().WithLevel(DEBUG).WithFormatter(JSONFormatter{}).AddWriter(ConsoleWriter{}).AddWriter(fileWriter).Build()

	defer log.Close()

	log.Info("service started", Feild{"port", 8080})
	reqLog := log.With(Feild{"request_id", "abc-123"})
	reqLog.Debug("handling request")
	reqLog.Error("db timeout", Feild{"retry", true})
}
