package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	BaseURL *string
	Port *int
	httpClient *http.Client
	St *Start
	OnCall *OnCalls
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{httpClient: httpClient}
	c.BaseURL = flag.String("base_url","http://localhost", "Server base url.")
	c.Port = flag.Int("port", 8000, "Server port")
	flag.Parse()
	c.St = &Start{}
	c.OnCall = &OnCalls{}
	return c
}
func (c *Client) ReqStartAPI(user_key string, problem_id int, number_of_elevators int) {
	flag.Parse()
	url := *c.BaseURL + ":" +
		strconv.Itoa(*c.Port) + "/" + "start" + "/" +
		user_key + "/" +
		strconv.Itoa(problem_id) + "/" +
		strconv.Itoa(number_of_elevators)

	resp, err := c.httpClient.Post(url,"application/json",
		nil)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
	if err := json.Unmarshal(respBody, c.St) ; err != nil {
		panic(err)
	}
}

func (c *Client) ReqOnCallAPI (){
	url := *c.BaseURL + ":" + strconv.Itoa(*c.Port) + "/" + "oncalls"
	nReq, _ := http.NewRequest("GET", url, nil)
	nReq.Header.Add("X-Auth-Token",c.St.Token)
	resp, err := c.httpClient.Do(nReq)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(respBody, c.OnCall) ; err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))

}
func main() {
	c := NewClient(&http.Client{})
	c.ReqStartAPI("tester",0,4)
	c.ReqOnCallAPI()
}