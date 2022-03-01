package dynamics

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func getQADynamicAPI(pn int, ps int, log *zap.Logger) (string, error){
	url := fmt.Sprintf("https://api.bilibili.com/x/space/dynamic/search?keyword=%E6%AF%8F%E5%91%A8QA&mid=703007996&pn=%d&ps=%d&platform=web",pn ,ps)
	method := "GET"

	client := &http.Client {
	}
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
	fmt.Println(string(body))
}