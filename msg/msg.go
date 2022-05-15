package msg

type Msg struct {
	Action  string
	Content interface{}
}

func New(action string, content interface{}) *Msg {
	return &Msg{action, content}
}
