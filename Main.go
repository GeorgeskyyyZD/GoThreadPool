package main

import (
	"./DaoLib"
	"fmt"
)

func main() {
	threadPool := DaoLib.CreateThreadPool(300)
	for i := 0; i < 1000; i++ {
		var thread DaoLib.Thread = DaoLib.Thread{ThreadCode:i, IRunnable:&MyThread{}}
		threadPool.AddTask(thread)
	}
	threadPool.ExecuteAllTasks()
	fmt.Println("execute over")
}

type MyThread struct {
	DaoLib.Thread
}

func (this *MyThread)DoingBackground() bool {
	fmt.Println("my DoingBackground")
	return true
}

func (this *MyThread)PostExecute() bool {
	fmt.Println("my PostExecute")
	return false
}
