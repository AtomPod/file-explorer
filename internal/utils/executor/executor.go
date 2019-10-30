package executor

//ExecuteFunc 执行函数
type ExecuteFunc func() error

//Action 事件
type Action interface {
	Tag() string    //标签
	Execute() error //执行函数
}

//Executor 执行器
type Executor struct {
	actionList []Action
}

//Add 添加一个执行事件
func (e *Executor) Add(a Action) {
	e.actionList = append(e.actionList, a)
}

//AddFuncWithTag 添加一个执行函数
func (e *Executor) AddFuncWithTag(tag string, fn func() error) {
	e.Add(&fnAction{tag: tag, fn: fn})
}

//AddFunc 添加一个执行函数
func (e *Executor) AddFunc(fn func() error) {
	e.AddFuncWithTag("", fn)
}

//Execute 执行所有操作
func (e *Executor) Execute(stopOnErr bool, errCallback func(a Action, err error)) error {
	for _, action := range e.actionList {
		err := action.Execute()
		if err != nil {
			if stopOnErr {
				return err
			}
			if errCallback != nil {
				errCallback(action, err)
			}
		}
	}
	return nil
}
