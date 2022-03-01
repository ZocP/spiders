package main

import (
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/dynamics"
)

func main() {
	log, _ := zap.NewDevelopment()
	dynamics.GetDynamicsIDs(log, config.Config{})
}
