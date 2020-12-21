package goo

func AsyncFunc(fn func()) {
	go func(fn func()) {
		defer func() {
			if err := recover(); err != nil {
				Log.Trace().Error(err)
			}
		}()
		fn()
	}(fn)
}
