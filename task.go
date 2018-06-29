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
	StateOverdue   = "overdue"
)

var (
	taskStateMutex sync.Mutex
	taskPool       sync.Pool
	taskChan     = make(chan *Task, runtime.NumCPU())
	taskState    = make(map[string]string)
)

type Task struct {
	Param   map[string]interface{}
	Factory []FacFunc
	UUID    string
	Expiration int64
}

type FacFunc func(string, map[string]interface{}) (string, error)

func NewTask(param  map[string]interface{}, factory []FacFunc, d time.Duration) *Task {

	var expiration int64
	if d > 0 {
		expiration = time.Now().Add(d).UnixNano()
	} else {
		expiration = -1
	}

	t := taskPool.Get()
	if t == nil {
		return &Task{
			param,
			factory,
			getUUID(20),
			expiration,
		}
	} else {
		task := t.(*Task)
		(*task).Param = param
		(*task).Factory = factory
		(*task).UUID = getUUID(20)
		(*task).Expiration = expiration
		return task
	}

}

func AddTask(task *Task) string {
	go func() {
		taskChan <- task
	}()
	uuid := (*task).UUID
	UpdateTaskState(uuid, StateWaiting)
	return uuid
}

func UpdateTaskState(uuid, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	taskState[uuid] = state
}

func GetTaskState(uuid string) (state string) {
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
	var taskUUID string
	var err error
	for {
		task := <-taskChan

		if ((*task).Expiration > 0 && time.Now().UnixNano() < (*task).Expiration) || (*task).Expiration < 0 {
			for _, f := range (*task).Factory {
				taskUUID, err = f((*task).UUID, (*task).Param)
			}
			if err != nil {
				UpdateTaskState(taskUUID, StateError)
				taskPool.Put(task)
			} else {
				UpdateTaskState(taskUUID, StateCompleted)
				taskPool.Put(task)
			}
		} else {
			UpdateTaskState(taskUUID, StateOverdue)
			taskPool.Put(task)
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