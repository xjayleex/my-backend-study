package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)


var (
	baseURL = flag.String("base_url","http://localhost", "Server base url.")
	port = flag.Int("port", 8000, "Server port")
)

func start(user_key string, problem_id int, number_of_elevators int) {
	var start Start
	flag.Parse()
	url := *baseURL + ":" +
		strconv.Itoa(*port) + "/" + "start" + "/" +
		user_key + "/" +
		strconv.Itoa(problem_id) + "/" +
		strconv.Itoa(number_of_elevators)
	fmt.Println("Url is " +url)
	client := &http.Client{}
	resp, err := client.Post(url,"application/json",
		nil)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
	if err := json.Unmarshal(respBody, &start) ; err != nil {
		panic(err)
	}
}

func main() {
	start("tester",0,4)
}