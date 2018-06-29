package task

import (
	"sync"
	"runtime"
	"time"
	"math/rand"
)

const (
	StateWaiting   = "waiting"
	StateCompleted = "completed"
	StateError     = "failed"
	StateNone      = "none"
)

var taskStateMutex sync.Mutex

var taskChan = make(chan Task, runtime.NumCPU())

var taskState = make(map[string]string)

type Task struct {
	Param   map[string]interface{}
	Factory []FacFunc
	UUID    string
}

type FacFunc func(string, map[string]interface{}) (string, error)

func NewTask(param  map[string]interface{}, factory []FacFunc) Task {
	return Task{
		param,
		factory,
		getUUID(20),
	}
}

func AddTask(task Task) string {
	go func() {
		taskChan <- task
	}()
	uuid := task.UUID
	UpdateTaskState(uuid, StateWaiting)
	return uuid
}

func UpdateTaskState(uuid, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	taskState[uuid] = state
}

func LoadTaskState(uuid string) (state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[uuid]
	if !exists {
		state = StateNone
	} else {
		state = resultState
	}
	return
}

func taskReceiver() {
	for {
		task := <-taskChan
		var curTaskHash string
		var err error
		for _, f := range task.Factory {
			curTaskHash, err = f(task.UUID, task.Param)
		}
		if err != nil {
			UpdateTaskState(curTaskHash, StateError)
		} else {
			UpdateTaskState(curTaskHash, StateCompleted)
		}
	}
}

func InitTaskReceiver(num int) {
	for i := 0; i < num; i++ {
		go taskReceiver()
	}
}

func getUUID(length int64) string {
	ele := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "v", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ele, _ = random(ele)
	uuid := ""
	var i int64
	for i = 0; i < length; i++ {
		rand.Seed(time.Now().Unix() + i)
		uuid += ele[rand.Intn(59)]
	}
	return uuid
}

func random(strings []string) ([]string, error) {
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := make([]string, 0)
	for i := 0; i < len(strings); i++ {
		str = append(str, strings[i])
	}
	return str, nil
}