package dynamics

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getQADynamicAPI(pn int, ps int, log *zap.Logger) (string, error) {
	url := "https://api.bilibili.com/x/space/dynamic/search?keyword=%E6%AF%8F%E5%91%A8QA&mid=703007996&pn=" + strconv.Itoa(pn) + "+&ps=" + strconv.Itoa(ps) + "&platform=web"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil

}

func getArticle(id string, log *zap.Logger) (io.ReadCloser, error) {
	url := "https://www.bilibili.com/read/cv" + id
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Error("new request", zap.Error(err))
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error("do", zap.Error(err))
		return nil, err
	}
	return res.Body, nil
}
