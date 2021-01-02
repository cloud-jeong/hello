package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
)

type Kube interface {
	DoSome(arg string) string
}

// Data is data
type Data struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
	Etc1 string `json:"Etc1,omitempty"`
	Etc2 string `json:"Etc2,omitempty"`
	Temp string `json:"-"`
}

func (d *Data) DoSome(arg string) string {
	return arg
}

func (d *Data) DoSome1(arg string) string {
	return arg
}

func (d *Data) init() {
	d.Name = "cloud.jung"
	d.Age = 30
	d.Etc1 = ""
	d.Etc2 = ""
	d.Temp = ""
}

func NewData(name string) *Data {
	data := &Data{
		Name: name,
		Age:  10,
		Etc1: "",
		Etc2: "",
		Temp: "",
	}

	//data.init()
	return data
}

func (d *Data) GetInstance() *Data {
	return nil
}

func main() {
	var aa Kube
	aa = NewData("cloud.jung@acronsoft.io")
	klog.Infoln(aa.DoSome("aaaaaa"))

	klog.Infoln("[version] retrieving version info")

	d := NewData("cloud.jung@acronsoft.io")
	klog.Infoln(d.Name)

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
