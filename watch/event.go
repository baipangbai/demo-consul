package watch

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/consul/api/watch"
)

func Event() {
	spew.Dump("watch event start")
	var (
		err    error
		params map[string]interface{}
		plan   *watch.Plan
		ch     chan int
	)
	ch = make(chan int, 1)

	params = make(map[string]interface{})
	params["type"] = "event"
	params["name"] = "event-test"

	plan, err = watch.Parse(params)
	if err != nil {
		panic(err)
	}
	plan.Handler = func(index uint64, result interface{}) {
		spew.Dump("index", index)
		spew.Dump("plan.Handler is ", result)
		ch <- 1
	}

	go func() {
		// your consul agent addr
		if err = plan.Run("127.0.0.1:8500"); err != nil {
			panic(err)
		}
	}()

	for {
		<-ch
		spew.Dump("event-test get changed")
	}
}
