# go-task

simple go asynchronous task generator tool

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
	"testing"
)

func main() {

	// init
	TaskReceiveInit(runtime.NumCPU())

	// add task
	AddTask(NewTask(
		map[string]interface{}{},  // parameter
		[]makeFunc{func(uuid string, param map[string]interface{}) (string, error) {
			fmt.Println(uuid)
			fmt.Println(param)
			return "ok", nil
		}}),
	)

	time.Sleep(time.Second * 5)
}
```