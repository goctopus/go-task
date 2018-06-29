package task

import (
	"runtime"
	"fmt"
	"time"
	"testing"
)

func testBuild(t *testing.T) {

	// init
	InitTaskReceiver(runtime.NumCPU())

	// add task
	AddTask(NewTask(
		map[string]interface{}{
			"paramA" : "value",
		},  // parameter
		[]FacFunc{func(uuid string, param map[string]interface{}) (string, error) {
			fmt.Println(uuid)
			fmt.Println(param)
			return "ok", nil
		}}),
	)

	time.Sleep(time.Second * 5)
}