package goo

func AsyncFunc(fn func()) {
	go func(fn func()) {
		defer func() {
			if err := recover(); err != nil {
				Log.Error(err)
			}
		}()
		fn()
	}(fn)
}
