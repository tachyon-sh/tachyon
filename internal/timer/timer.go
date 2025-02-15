package timer

import (
	"fmt"
	"time"
)

type TaskTimer struct {
	startTime time.Time
	taskName  string
}

func StartTask(name string) *TaskTimer {
	fmt.Printf("⏳ %s...\n", name)
	return &TaskTimer{startTime: time.Now(), taskName: name}
}

func (t *TaskTimer) Stop() {
	duration := time.Since(t.startTime)
	fmt.Printf("✅ %s завершено за %.2f сек\n", t.taskName, duration.Seconds())
}