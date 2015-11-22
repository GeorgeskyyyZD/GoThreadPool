package DaoLib
import (
	"fmt"
	"sync"
)

type ThreadPool struct {
	Signals []chan int
	Tasks   map[int]Thread
}

var pool ThreadPool
var totalTask int
var waitingSignal chan int
var selfWG sync.WaitGroup

func CreateThreadPool(capacity int) ThreadPool {
	waitingSignal = make(chan int)
	totalTask = capacity
	pool = ThreadPool{Tasks: make(map[int]Thread)}
	return pool
}

func (pool *ThreadPool)AddTask(thread Thread) {
	pool.Tasks[thread.ThreadCode] = thread
	pool.Signals = make([]chan int, len(pool.Tasks))
}

func (pool *ThreadPool)ExecuteAllTasks() {
	selfWG.Add(1)
	go func() {
		taskRunningCount := 0
		taskFinishCount := 0
		for threadCode, singleTask := range pool.Tasks {
			fmt.Println("current index:", taskRunningCount, ",address:", &singleTask)
			pool.Signals[taskRunningCount] = make(chan int)
			go pool.executeSingleTask(singleTask, pool.Signals[taskRunningCount], func() {
				delete(pool.Tasks, threadCode)
				taskFinishCount++
				if taskFinishCount == taskRunningCount {
					waitingSignal <- 1
				}
			})
			taskRunningCount++;
			if taskRunningCount >= totalTask {
				<-waitingSignal
			}
		}
		for _, signal := range pool.Signals {
			<-signal
		}
		selfWG.Done()
	}()
	selfWG.Wait()
}

func (pool *ThreadPool)executeSingleTask(thread Thread, signal chan int, finish func()) bool {
	hasBackGroundTaskFinish := thread.DoingBackground()
	if hasBackGroundTaskFinish {
		thread.PostExecute()
	}
	finish()
	signal <- thread.ThreadCode
	return true
}
