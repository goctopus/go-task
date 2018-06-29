# go-task

simple asynchronous task generator tool writed by go

## install

```
go get github.com/chenhg5/go-task
```

## usage

```
import (
	"runtime"
	"fmt"
	"time"
	"github.com/chenhg5/go-task"
)

func main() {

	// init
	task.InitTaskReceiver(runtime.NumCPU())

	// add task: parameter, taskList, expiration
	task.AddTask(task.NewTask(
		map[string]interface{}{
            "paramA" : "value",
        },  // parameter
		[]task.FacFunc{func(uuid string, param map[string]interface{}) (string, error) {
			fmt.Println(uuid)
			fmt.Println(param)
			return "ok", nil
		}}, -1),
	)

	time.Sleep(time.Second * 5)
}
```

## todo

- [ ] get the number of statistical state
- [ ] api for putting an end to a task
- [ ] api for clearing task list
- [ ] add task delay list