package main

import (
	"bytes"
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
	fmt.Println(string(respBody))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(respBody, c.St) ; err != nil {
		panic(err)
	}
	fmt.Println("Requested Start API ... / Status -> " + resp.Status)

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
	fmt.Println(string(respBody))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(respBody, c.OnCall) ; err != nil {
		panic(err)
	}
	fmt.Println("Requested Oncall ... / Status -> " + resp.Status)
}

func (c *Client) ActionAPI (cmd *Commands) {
	result := &ActionResult{}
	url := *c.BaseURL + ":" + strconv.Itoa(*c.Port) + "/" + "action"
	b, err := json.Marshal(*cmd)
	fmt.Println(*cmd)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(b)
	if err != nil {
		panic(err)
	}

	nReq, _ := http.NewRequest("POST", url, buf)
	nReq.Header.Add("X-Auth-Token",c.St.Token)
	nReq.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(nReq)
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(respBody, result) ; err != nil {
		panic(err)
	}
	fmt.Println("Requested Action ... / Status -> " + resp.Status)
}
