package threading

import (
	"github.com/51xpx/pkgx/pkg/lang"
	"github.com/51xpx/pkgx/pkg/rescue"
)

type TaskRunner struct {
	limitChan chan lang.PlaceholderType
}

func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan lang.PlaceholderType, concurrency),
	}
}

func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- lang.Placeholder

	// threading.GoSafe(func() {
	GoSafe(func() {
		defer rescue.Recover(func() {
			<-rp.limitChan
		})

		task()
	})
}
