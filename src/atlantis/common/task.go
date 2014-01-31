package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const taskIDSize = 20

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
	ID       string
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

func createTaskID() string {
	return CreateRandomID(taskIDSize)
}

func (t *TaskTracker) ListIDs(types []string) []string {
	typesMap = make(map[string]bool, len(types))
	for _, typ := range types {
		typesMap[typ] = true
	}
	ids := []string{}
	t.Lock()
	for id, task := range t.tasks {
		if typesMap[task.Name] {
			ids = append(ids, id)
		}
	}
	t.Unlock()
	return ids
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

func (t *TaskTracker) ReserveTaskID(task *Task) string {
	t.Lock()
	requestID := createTaskID()
	for _, present := t.tasks[requestID]; present; _, present = t.tasks[requestID] {
		requestID = createTaskID()
	}
	t.tasks[requestID] = task // reserve request id
	t.Unlock()
	task.Lock()
	task.ID = requestID
	task.Unlock()
	return requestID
}

func (t *TaskTracker) ReleaseTaskID(id string) {
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
	t.Tracker.ReserveTaskID(t)
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
	t.Tracker.ReserveTaskID(t)
	t.RLock()
	r.ID = t.ID
	t.RUnlock()
	go func() error {
		t.Log("Begin %s", t.Description)
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
			t.Tracker.ReleaseTaskID(t.ID)
		})
	} else {
		t.Tracker.ReleaseTaskID(t.ID)
	}
	return err
}

func (t *Task) Log(format string, args ...interface{}) {
	t.RLock()
	log.Printf("[RPC]["+t.Name+"]["+t.ID+"] "+format, args...)
	t.RUnlock()
}

func (t *Task) LogStatus(format string, args ...interface{}) {
	t.Log(format, args...)
	t.Lock()
	t.StatusTime = time.Now()
	t.Status = fmt.Sprintf(format, args...)
	t.Unlock()
}

func (t *Task) AddWarning(warn string) {
	t.Lock()
	if t.Warnings == nil {
		t.Warnings = []string{warn}
	} else {
		t.Warnings = append(t.Warnings, warn)
	}
	t.Unlock()
	t.Log("WARNING: %s", warn)
}

type TaskStatus struct {
	Name        string
	Description string
	Status      string
	Warnings    []string
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
		"Warnings":    t.Warnings,
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
Warnings    : %v
Done        : %t
StartTime   : %s
StatusTime  : %s
EndTime     : %s`, t.Name, t.Description, t.Status, t.Warnings, t.Done, t.StartTime, t.StatusTime,
		t.EndTime)
}

func (t *TaskStatus) CopyTaskStatus() *TaskStatus {
	return &TaskStatus{t.Name, t.Description, t.Status, t.Warnings, t.Done, t.StartTime, t.StatusTime,
		t.EndTime}
}
