package singleflight

import sf "golang.org/x/sync/singleflight"

var SfGroup *sf.Group

func SingleFlightInit() {
	SfGroup = new(sf.Group)
}

func GetSingleFlight() *sf.Group {
	return SfGroup
}

func Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return SfGroup.Do(key, fn)
}

//TODO async timeout-context
func DoChan(key string, fn func() (any, error)) <-chan sf.Result {
	return SfGroup.DoChan(key, fn)
}

func Forget(key string) {
	SfGroup.Forget(key)
}
