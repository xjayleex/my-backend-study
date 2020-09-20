package main

import (
	"fmt"
	"net/http"
	"strconv"
)
// 사용할 엘리베이터의 수를 1대에서 4대까지 선택할 수 있다.
// 건물의 가장 낮은 층은 항상 1층이다.
// Call 목록은 건물 마다 다르게 주어진다.
// 유효하지 않은 Call은 없다.(존재하지 않는 층으로의 이동, 같은 층으로의 이동 등등)
// Token의 유효시간인 10분 안에 모든 Call을 처리하지 못하면 점수를 받을 수 없다.

// Problem 0
// 출근을 위해 집을 나선 라이언. 라이언이 살고있는 어피치 맨션은 총 5층 높이의 작은 맨션이다.
// 이동이 많지 않은 기본적인 엘리베이터 동작을 구현하여 승객을 수송해보자!
// 조건
// 엘리베이터의 최대 수용인원(=Call) : 8명
// 건물의 최고층 : 5층
// Call 수 : 6개

func solution0 (numElevators int) {
	c := NewClient(&http.Client{})
	c.ReqStartAPI("tester",0,numElevators)
	for _, e := range c.St.Elevators {
		fmt.Println("Elevator " + strconv.Itoa(e.Id) + " : " + e.Status)
	}
	c.ReqOnCallAPI()

	/*
	for !c.St.IsEnd {
		c.ReqOnCallAPI()
		for i := 0 ; i < numElevators ; i++ {
			switch c.St.Elevators[i].Status {
			case "STOPPED" :
				// Do something : STOP, UP, DOWN, OPEN
			case "UPWARD" :
				// Do something : STOP, UP
			case "DOWNWARD" :
				// Do something : STOP, DOWN
			case "OPENED" :
				// Do something : OPEN, CLOSE, ENTER, EXIT
			}
		}
	}*/

}
func main() {
	solution0(3)
	/*cmds := &Commands{
		Cmds:  &[]Command{
			{
				ElevatorID: 0, Command: "ENTER", CallIDs: &[]int{0},
			}, {
				ElevatorID: 1, Command: "STOP",
			},
		},
	}

	c.ActionAPI(cmds)*/
}