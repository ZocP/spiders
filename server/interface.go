package server

import "qa_spider/config"

type Server interface {
	Run() error
	Stop()
	ReloadConfig(config *config.Config) error
}
