package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const taskIdSize = 20

var (
	Tracker           = &TaskTracker{tasks: map[string]*Task{}}
	TaskStatusUnknown = &TaskStatus{Status: StatusUnknown}
)

func MaintenanceChecker(file string, interval time.Duration) {
	go Tracker.MaintenanceChecker(file, interval)
}

func NewTask(name string, executor TaskExecutor) *Task {
	task := &Task{Tracker: Tracker, Executor: executor}
	task.Status = StatusInit
	task.StatusTime = time.Now()
	task.Name = name
	task.Description = executor.Description()
	task.Request = executor.Request()
	return task
}

type TaskTracker struct {
	sync.RWMutex
	ResultDuration time.Duration
	Maintenance    bool
	tasks          map[string]*Task
}

type Task struct {
	sync.RWMutex
	TaskStatus
	Err      error
	Tracker  *TaskTracker
	Id       string
	Executor TaskExecutor
	Request  interface{}
	Result   interface{}
}

type TaskExecutor interface {
	Request() interface{}
	Result() interface{}
	Description() string
	Execute(t *Task) error
	Authorize() error
}

type TaskMaintenanceExecutor interface {
	AllowDuringMaintenance() bool
}

func createTaskId() string {
	return CreateRandomId(taskIdSize)
}

func (t *TaskTracker) SetMaintenance(on bool) {
	t.Lock()
	t.Maintenance = on
	t.Unlock()
}

func (t *TaskTracker) UnderMaintenance() bool {
	t.RLock()
	maint := t.Maintenance
	t.RUnlock()
	return maint
}

func (t *TaskTracker) MaintenanceChecker(file string, interval time.Duration) {
	for {
		if _, err := os.Stat(file); err == nil {
			// maintenance file exists
			if !t.UnderMaintenance() {
				log.Println("Begin Maintenance")
				t.SetMaintenance(true)
			}
		} else {
			// maintenance file doesn't exist or there is an error looking for it
			if t.UnderMaintenance() {
				log.Println("End Maintenance")
				t.SetMaintenance(false)
			}
		}
		time.Sleep(interval)
	}
}

func (t *TaskTracker) Idle(checkTask *Task) bool {
	idle := true
	t.RLock()
	for _, task := range t.tasks {
		if task != checkTask && !task.Done {
			idle = false
			break
		}
	}
	t.RUnlock()
	return idle
}

func (t *TaskTracker) ReserveTaskId(task *Task) string {
	t.Lock()
	requestId := createTaskId()
	for _, present := t.tasks[requestId]; present; _, present = t.tasks[requestId] {
		requestId = createTaskId()
	}
	t.tasks[requestId] = task // reserve request id
	t.Unlock()
	task.Lock()
	task.Id = requestId
	task.Unlock()
	return requestId
}

func (t *TaskTracker) ReleaseTaskId(id string) {
	t.Lock()
	delete(t.tasks, id)
	t.Unlock()
}

func (t *TaskTracker) Status(id string) (*TaskStatus, error) {
	t.RLock()
	task := t.tasks[id]
	t.RUnlock()
	if task != nil {
		task.RLock()
		status := task.CopyTaskStatus()
		err := task.Err
		task.RUnlock()
		return status, err
	}
	return TaskStatusUnknown, errors.New("Unknown Task Status")
}

func (t *TaskTracker) Result(id string) interface{} {
	t.RLock()
	task := t.tasks[id]
	t.RUnlock()
	if task != nil {
		task.RLock()
		result := task.Result
		task.RUnlock()
		return result
	}
	return nil
}

func (t *Task) Authorize() error {
	return t.Executor.Authorize()
}

func (t *Task) Run() error {
	if t.Tracker.UnderMaintenance() {
		executor, ok := t.Executor.(TaskMaintenanceExecutor)
		if !ok || !executor.AllowDuringMaintenance() {
			return t.End(errors.New("Under Maintenance"), false)
		}
	}
	t.Tracker.ReserveTaskId(t)
	t.Log("Begin %s", t.Description)
	t.Lock()
	t.StartTime = time.Now()
	t.Unlock()
	err := t.Executor.Authorize()
	if err != nil {
		return t.End(err, false)
	}
	return t.End(t.Executor.Execute(t), false)
}

func (t *Task) RunAsync(r *AsyncReply) error {
	if t.Tracker.UnderMaintenance() {
		executor, ok := t.Executor.(TaskMaintenanceExecutor)
		if !ok || !executor.AllowDuringMaintenance() {
			return t.End(errors.New("Under Maintenance"), false)
		}
	}
	t.Tracker.ReserveTaskId(t)
	t.RLock()
	r.Id = t.Id
	t.RUnlock()
	go func() error {
		t.Log("Begin")
		t.Lock()
		t.StartTime = time.Now()
		t.Unlock()
		err := t.Executor.Authorize()
		if err != nil {
			return t.End(err, true)
		}
		t.End(t.Executor.Execute(t), true)
		return nil
	}()
	return nil
}

func (t *Task) End(err error, async bool) error {
	logString := fmt.Sprintf("End %s", t.Description)
	t.Lock()
	t.Result = t.Executor.Result()
	t.EndTime = time.Now()
	t.StatusTime = t.EndTime
	if err == nil {
		t.Status = StatusDone
		t.Done = true
	} else {
		t.Status = StatusError
		t.Err = err
		t.Done = true
		logString += fmt.Sprintf(" - Error: %s", err.Error())
	}
	t.Unlock()
	t.Log(logString)
	if async {
		time.AfterFunc(t.Tracker.ResultDuration, func() {
			// keep result around for 30 min in case someone wants to check on it
			t.Tracker.ReleaseTaskId(t.Id)
		})
	} else {
		t.Tracker.ReleaseTaskId(t.Id)
	}
	return err
}

func (t *Task) Log(format string, args ...interface{}) {
	t.RLock()
	log.Printf("[RPC]["+t.Name+"]["+t.Id+"] "+format, args...)
	t.RUnlock()
}

func (t *Task) LogStatus(format string, args ...interface{}) {
	t.Log(format, args...)
	t.Lock()
	t.StatusTime = time.Now()
	t.Status = fmt.Sprintf(format, args...)
	t.Unlock()
}

type TaskStatus struct {
	Name        string
	Description string
	Status      string
	Done        bool
	StartTime   time.Time
	StatusTime  time.Time
	EndTime     time.Time
}

func (t *TaskStatus) Map() map[string]interface{} {
	return map[string]interface{}{
		"Name":        t.Name,
		"Description": t.Description,
		"Status":      t.Status,
		"Done":        t.Done,
		"StartTime":   t.StartTime,
		"StatusTime":  t.StatusTime,
		"EndTime":     t.EndTime,
	}
}

func (t *TaskStatus) String() string {
	return fmt.Sprintf(`%s
Description : %s
Status      : %s
Done        : %t
StartTime   : %s
StatusTime  : %s
EndTime     : %s`, t.Name, t.Description, t.Status, t.Done, t.StartTime, t.StatusTime,
		t.EndTime)
}

func (t *TaskStatus) CopyTaskStatus() *TaskStatus {
	return &TaskStatus{t.Name, t.Description, t.Status, t.Done, t.StartTime, t.StatusTime,
		t.EndTime}
}
