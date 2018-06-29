package task

import (
	"runtime"
	"fmt"
	"time"
	"testing"
)

func testBuild(t *testing.T) {

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