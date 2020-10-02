package main

type Common struct {
	Limit int
}

type Call struct {
	Id int `json:"id"` //Call 고유 번호
	Timestamp int `json:"timestamp"` // 해당 Call이 발생한 timestamp
	Start int `json:"start"` // 출발 층
	End int `json:"end"` // 가려는 층
}
type OnCalls struct {
	Token string `json:"token"`
	Timestamp int `json:"timestamp"`
	Elevators []Elevator `json:"elevators"`
	Calls []Call `json:"calls"`
	IsEnd bool `json:"is_end"`
}

type Start struct {
	Token string `json:"token"`
	Timestamp int `json:"timestamp"`
	Elevators []Elevator `json:"elevators"`
	IsEnd bool `json:"is_end"`
}
type Commands struct {
	Cmds *[]Command `json:"commands"`
}
type Command struct {
	ElevatorID int `json:"elevator_id"`
	Command string `json:"command"`
	CallIDs *[]int `json:"calls_ids,omitempty"`
}

type ActionResult struct {
	Token string
	Timestamp int
	Elevator []Elevator
	IsEnd bool `json:"is_end"`
}
// 엘리베이터는 여러 대가 존재하며 모두 사용할 수도 있고 일부만 사용해도 된다.
// 엘리베이터에 명령을 내려 각각의 엘리베이터를 층을 이동하거나 멈추고, 문을 열거나 닫고, 승객을 태우거나 내려 줄 수 있다.
// 엘리베이터는 정원이 있어 정해진 수 이상의 승객을 태울 수 없다.

type Elevator struct {
	Id int `json:"id"`
	Floor int `json:"floor"`
	Passengers []Call `json:"passengers"`
	Status string `json:"status"`

}


// 엘리베이터에는 현재 상태를 표현하는 status가 있으며, 값으로는 STOPPED, OPENED, UPWARD, DOWNWARD가 있다.
// 사용할 수 있는 명령은 다음과 같다.
// STOP		엘리베이터를 멈춘다. 현재 층에 머무르기 원하는 경우 STOP 명령을 통해 머무를 수 있다.
// UP		엘리베이터를 한 층 올린다. 최상층인 경우 현재 층을 유지한다.
// DOWN		엘리베이터를 한 층 내린다. 1층인 경우 1층을 유지한다.
// OPEN		엘리베이터의 문을 연다. 엘리베이터의 문이 열린 상태를 유지하기 위해서는 OPEN 명령을 사용한다.
// CLOSE	엘리베이터의 문을 닫는다.
// ENTER	엘리베이터에 승객을 태운다.
// EXIT		엘리베이터의 승객을 내린다. 목적지가 아닌 곳에서 내린 경우, OnCall
//			API의 calls에 내린 층과 내린 시점의 timestamp로 변경되어 다시 들어가게 된다.
