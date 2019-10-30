package executor

type fnAction struct {
	tag string
	fn  func() error
}

func (a *fnAction) Tag() string {
	return a.tag
}

func (a *fnAction) Execute() error {
	return a.fn()
}
