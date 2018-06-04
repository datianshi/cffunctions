package cffunctions

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Org struct {
	Api API
	ORG cfclient.Org
}

func (org Org) GetAPI() API {
	return org.Api
}

func (org Org) GetName() string {
	return org.ORG.Name
}

func (org Org) GetGuid() string {
	return org.ORG.Guid
}

func (org Org) GetParent() (CFObject, error) {
	return nil, nil
}

type Space struct {
	Api   API
	SPACE cfclient.Space
}

func (space Space) GetAPI() API {
	return space.Api
}

func (space Space) GetName() string {
	return space.SPACE.Name
}

func (space Space) GetGuid() string {
	return space.SPACE.Guid
}

func (space Space) GetParent() (CFObject, error) {
	org, err := space.SPACE.Org()
	if err != nil {
		return nil, err
	}
	return Org{
		Api: space.GetAPI(),
		ORG: org,
	}, nil
}

type OrgQuota struct {
	Api   API
	QUOTA cfclient.OrgQuota
}

type SpaceQuota struct {
	Api   API
	QUOTA cfclient.OrgQuota
}

type CFObject interface {
	GetAPI() API
	GetName() string
	GetGuid() string
	GetParent() (CFObject, error)
}

type Object struct {
	Retriever   func(API) ([]interface{}, error)
	Constructor func(API, interface{}) CFObject
}

var (
	OrgFunc ObjectFunc = func(api API) ([]CFObject, error) {
		objects, err := api.ListOrgs()
		if err != nil {
			return nil, err
		}
		retObjects := make([]CFObject, 0)
		for _, object := range objects {
			o := Org{
				Api: api,
				ORG: object,
			}
			retObjects = append(retObjects, o)
		}
		return retObjects, nil
	}

	SpaceFunc ObjectFunc = func(api API) ([]CFObject, error) {
		objects, err := api.ListSpaces()
		if err != nil {
			return nil, err
		}
		retObjects := make([]CFObject, 0)
		for _, object := range objects {
			o := Space{
				Api:   api,
				SPACE: object,
			}
			retObjects = append(retObjects, o)
		}
		return retObjects, nil
	}
)
