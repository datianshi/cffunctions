package cffunctions

import (
	"errors"
	"fmt"
)

const (
	//ConcurrencyCapacity control the concurrency
	ConcurrencyCapacity int = 3
)

//CloudFoundry CloudFoundry
type CloudFoundry struct {
	Api API
}

type Looper func() ([]CFObject, error)
type Action func(CFObject) error
type Filter func(CFObject) bool

type ObjectFunc func(API) ([]CFObject, error)

//Each EachObject
func (cf *CloudFoundry) Each(f ObjectFunc) Looper {
	return func() ([]CFObject, error) {
		return f(cf.Api)
	}
}

func (looper Looper) Action(action Action) Looper {
	return func() ([]CFObject, error) {
		objects, err := looper()
		if err != nil {
			return nil, err
		}
		for _, o := range objects {
			err = action(o)
			if err != nil {
				return nil, err
			}
		}
		return objects, err
	}
}

func (looper Looper) Parallel(action Action) Looper {
	return func() ([]CFObject, error) {
		errorAggregator := make(chan error)
		semaphone := make(chan bool, ConcurrencyCapacity)
		objects, err := looper()
		if err != nil {
			return nil, err
		}
		for _, object := range objects {
			semaphone <- true
			go func(o CFObject) {
				err = action(o)
				<-semaphone
				errorAggregator <- err
			}(object)
		}
		var retErrorString string
		for i := 0; i < cap(semaphone); i++ {
			semaphone <- true
		}
		for i := 0; i < len(objects); i++ {
			r := <-errorAggregator
			if r != nil {
				retErrorString = fmt.Sprintf("%s\n%s", retErrorString, r)
			}
		}
		if retErrorString != "" {
			return objects, errors.New(retErrorString)
		}
		return objects, nil
	}
}

func (looper Looper) Filter(filter Filter) Looper {
	return func() ([]CFObject, error) {
		objects, err := looper()
		if err != nil {
			return nil, err
		}
		retObjects := make([]CFObject, 0)
		for _, object := range objects {
			if filter(object) {
				retObjects = append(retObjects, object)
			}
		}
		return retObjects, nil
	}
}

func MultipleActions(actions ...Action) Action {
	return func(object CFObject) error {
		for _, action := range actions {
			err := action(object)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
