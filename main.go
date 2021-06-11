package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/consul/api/watch"

	"github.com/hashicorp/consul/api"
)

//main 监测服务变化
func main() {
	// Create a Consul API client
	var (
		err    error
		params map[string]interface{}
		plan   *watch.Plan
		ch     chan int
	)
	ch = make(chan int, 1)

	params = make(map[string]interface{})
	params["type"] = "service"
	params["service"] = "test"
	params["passingonly"] = false
	params["tag"] = "SERVER"
	plan, err = watch.Parse(params)
	if err != nil {
		panic(err)
	}
	plan.Handler = func(index uint64, result interface{}) {
		if entries, ok := result.([]*api.ServiceEntry); ok {
			fmt.Printf("serviceEntries:%v", entries)
			// your code
			ch <- 1
		}
	}

	go func() {
		// your consul agent addr
		if err = plan.Run("127.0.0.1:8500"); err != nil {
			panic(err)
		}
	}()
	go http.ListenAndServe(":8080", nil)
	go register()
	for {
		<-ch
		fmt.Printf("get change")
	}
}

func register() {
	var (
		err    error
		client *api.Client
	)
	client, err = api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:   "",
		Name: "test",
		Tags: []string{"SERVER"},
		Port: 8080,
		Check: &api.AgentServiceCheck{
			HTTP: "",
		},
	})
	if err != nil {
		panic(err)
	}

}
