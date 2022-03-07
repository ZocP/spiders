package config

import (
	"bufio"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type Config struct {
	Server struct {
		Port      string `json:"port"`
		AllowCors bool   `json:"allow_cors"`
	} `json:"server"`
	Internal struct {
		QASpider struct {
			AutoUpdate bool   `json:"auto_update"`
			UpdateDate string `json:"update_date"`
			Writer     struct {
				Type     int `json:"writer_type"`
				LocalTxt struct {
					Path string `json:"file_path"`
				} `json:"local_txt"`
				LocalCSV struct {
				} `json:"local_csv"`
				RemoteDB struct {
					//	DB DB `json:"db,omitempty"`
				} `json:"remote_db"`
			} `json:"writer"`
		} `json:"qa_spider"`
	}
	Services struct {
		QueryQA struct {
			Longest  int `json:"longest"`
			Shortest int `json:"shortest"`
			//DB DB `json:"db,omitempty"`
		} `json:"query_qa"`
	} `json:"services"`
}

func InitConfig(log *zap.Logger) *Config {
	config, ok := readFromFiles(log)
	log.Info("config initializing")
	if !ok {
		log.Fatal("invalid config")
	}
	return config
}

func readFromFiles(log *zap.Logger) (*Config, bool) {
	config := &Config{}
	path := "./files/"
	log.Info("file path: " + path)
	if !pathExists(path) {
		if err := os.MkdirAll(path, 0777); err != nil {
			log.Error("making path: ", zap.Error(err))
		}
	}
	if !pathExists(path + "cfg.json") {
		makeBlankConfig(path+"cfg.json", log)
		log.Info("config file not found, generated new file, please fill in the config.")
		os.Exit(1)
		return nil, false
	}
	f, err := os.Open(path + "cfg.json")
	defer f.Close()
	if err != nil {

		log.Fatal("open file error")
	}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("reading from file: ", zap.Error(err))
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal("unmarshal from file: ", zap.Error(err))
		os.Exit(1)
	}
	return config, true
}

func makeBlankConfig(path string, log *zap.Logger) {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return
	}
	writer := bufio.NewWriter(f)
	var c Config
	setDefault(&c)
	raw, err := json.Marshal(c)
	if err != nil {
		log.Fatal("writing config: ", zap.Error(err))
	}
	_, err = writer.Write(raw)
	if err != nil {
		log.Fatal("writing config: ", zap.Error(err))
	}
	writer.Flush()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func setDefault(c *Config) {
	c.Internal.QASpider.Writer.LocalTxt.Path = "./files/spider/"
	c.Server.Port = ":8080"
	c.Server.AllowCors = true
	c.Internal.QASpider.AutoUpdate = true
	c.Internal.QASpider.Writer.Type = 0
	c.Internal.QASpider.UpdateDate = "WED"
	c.Services.QueryQA.Shortest = 3
	c.Services.QueryQA.Longest = 20

}
