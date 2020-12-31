package main

import (
	"encoding/json"
	"fmt"
)

// Data is data
type Data struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
	Etc1 string `json:"Etc1,omitempty"`
	Etc2 string `json:"Etc2,omitempty"`
	Temp string `json:"-"`
}

func main() {

	var jobj = Data{Name: "aaa", Age: 19, Etc1: "bbb", Etc2: "", Temp: "1123123"}
	jstrbyte, _ := json.Marshal(jobj)
	fmt.Println("---- jobj")
	fmt.Println(jobj)
	fmt.Println()
	fmt.Println("---- jstrbyte")
	fmt.Println(jstrbyte)
	fmt.Println()
	fmt.Println("---- string(jstrbyte)")
	fmt.Println(string(jstrbyte))
	fmt.Println()

	var jobj2 Data
	err := json.Unmarshal(jstrbyte, &jobj2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("---- jobj2")
	fmt.Println(jobj2)
	fmt.Println()
	fmt.Println("---- jobj2.Name, jobj2.Age, jobj2.Etc1, jobj2.Etc2")
	fmt.Println(jobj2.Name, jobj2.Age, jobj2.Etc1, jobj2.Etc2)
	fmt.Println()

}
