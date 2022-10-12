package faker

import (
	"math/rand"
	"time"
)

func MustFakeObject(obj interface{}, opts ...Option) interface{} {
	if err := FakeObject(obj, opts...); err != nil {
		panic(err)
	}
	return obj
}

func Name() string {
	var s struct {
		Name string
	}
	MustFakeObject(&s)
	return s.Name
}

func Word() string {
	return Name()
}

func RandomUnixTime() int64 {
	return rand.Int63n(time.Now().Unix())
}

func Timestamp() string {
	return time.Unix(RandomUnixTime(), 0).Format("2006-01-02 15:04:05")
}
