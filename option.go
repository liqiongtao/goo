package goo

type Option struct {
	Name  string
	Value interface{}
}

func (opt Option) Apply(m map[string]Option) {
	m[opt.Name] = opt
}

func (opt Option) String() string {
	if opt.Value == nil {
		return ""
	}
	return opt.Value.(string)
}

func (opt Option) Bool() bool {
	if opt.Value == nil {
		return false
	}
	return opt.Value.(bool)
}

func (opt Option) Int() int {
	if opt.Value == nil {
		return 0
	}
	return opt.Value.(int)
}

func (opt Option) Int64() int64 {
	if opt.Value == nil {
		return 0
	}
	return opt.Value.(int64)
}

func NewOption(name string, value interface{}) Option {
	return Option{Name: name, Value: value}
}
