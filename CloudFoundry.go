package cffunctions

import (
	"errors"
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

const (
	//ConcurrencyCapacity control the concurrency
	ConcurrencyCapacity int = 3
)

//CloudFoundry CloudFoundry
type CloudFoundry struct {
	Api API
}

type Org struct {
	Api API
	ORG cfclient.Org
}

type Space struct {
	Api   API
	SPACE cfclient.Space
}

type OrgAction func(Org) error
type SpaceAction func(Space) error
type OrgLooper func() ([]Org, error)
type SpaceLooper func() ([]Space, error)
type OrgFilter func(Org) bool
type SpaceFilter func(Space) bool

//EachOrg EachOrg
func (cf *CloudFoundry) EachOrg() OrgLooper {
	return func() ([]Org, error) {
		orgs, err := cf.Api.ListOrgs()
		if err != nil {
			return nil, err
		}
		retOrgs := make([]Org, 0)
		for _, org := range orgs {
			o := Org{Api: cf.Api, ORG: org}
			retOrgs = append(retOrgs, o)
		}
		return retOrgs, nil
	}
}

//EachSpace EachSpace
func (cf *CloudFoundry) EachSpace() SpaceLooper {
	return func() ([]Space, error) {
		spaces, err := cf.Api.ListSpaces()
		if err != nil {
			return nil, err
		}
		retSpaces := make([]Space, 0)
		for _, space := range spaces {
			s := Space{Api: cf.Api, SPACE: space}
			retSpaces = append(retSpaces, s)
		}
		return retSpaces, nil
	}
}

func (looper OrgLooper) Action(action OrgAction) OrgLooper {
	return func() ([]Org, error) {
		orgs, err := looper()
		if err != nil {
			return nil, err
		}
		for _, org := range orgs {
			err = action(org)
			if err != nil {
				return nil, err
			}
		}
		return orgs, err
	}
}

func (looper OrgLooper) Parallel(action OrgAction) OrgLooper {
	return func() ([]Org, error) {
		errorAggregator := make(chan error)
		semaphone := make(chan bool, ConcurrencyCapacity)
		orgs, err := looper()
		if err != nil {
			return nil, err
		}
		for _, org := range orgs {
			semaphone <- true
			go func(org Org) {
				err = action(org)
				<-semaphone
				errorAggregator <- err
			}(org)
		}
		var retErrorString string
		for i := 0; i < cap(semaphone); i++ {
			semaphone <- true
		}
		for i := 0; i < len(orgs); i++ {
			r := <-errorAggregator
			if r != nil {
				retErrorString = fmt.Sprintf("%s\n%s", retErrorString, r)
			}
		}
		if retErrorString != "" {
			return orgs, errors.New(retErrorString)
		}
		return orgs, nil
	}
}

func (looper SpaceLooper) Action(action SpaceAction) SpaceLooper {
	return func() ([]Space, error) {
		spaces, err := looper()
		if err != nil {
			return nil, err
		}
		for _, space := range spaces {
			err = action(space)
			if err != nil {
				return nil, err
			}
		}
		return spaces, err
	}
}

func (looper OrgLooper) Filter(filter OrgFilter) OrgLooper {
	return func() ([]Org, error) {
		orgs, err := looper()
		if err != nil {
			return nil, err
		}
		retOrgs := make([]Org, 0)
		for _, org := range orgs {
			if filter(org) {
				retOrgs = append(retOrgs, org)
			}
		}
		return retOrgs, nil
	}
}

func (looper SpaceLooper) Filter(filter SpaceFilter) SpaceLooper {
	return func() ([]Space, error) {
		spaces, err := looper()
		if err != nil {
			return nil, err
		}
		retSpaces := make([]Space, 0)
		for _, space := range spaces {
			if filter(space) {
				retSpaces = append(retSpaces, space)
			}
		}
		return retSpaces, nil
	}
}
