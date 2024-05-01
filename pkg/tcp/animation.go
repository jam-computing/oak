package tcp

import (
	"encoding/json"
	"fmt"
)

type Animation struct {
	Title          string        `json:"Title"`
	Artist         string        `json:"Artist"`
	Tick_Rate      uint16        `json:"Tick_Rate"`
	Frames         []interface{} `json:"Frames"`
	CollectionID   string        `json:"collectionId"`
	CollectionName string        `json:"collectionName"`
	Created        string        `json:"created"`
	Id             string        `json:"id"`
	Updated        string        `json:"updated"`
}

func MakeAnimation(data string) Animation {
	var a Animation
	err := json.Unmarshal([]byte(data), &a)

	if err != nil {
		fmt.Println(err)
		return Animation{}
	}

	return a
}
