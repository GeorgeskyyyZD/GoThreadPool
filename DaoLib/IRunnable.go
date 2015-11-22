package DaoLib

type IRunnable interface {
	DoingBackground() bool
	PostExecute() bool
}