package trigger

type Callback func()

type Trigger interface {
	AddCallback(callback Callback)
}
