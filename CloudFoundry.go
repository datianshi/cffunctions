package cffunctions

import (
	"io"
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//API Interface
type API interface {
	BindRoute(routeGUID, appGUID string) error
	CreateApp(appCreateRequest cfclient.AppCreateRequest) (app cfclient.App, err error)
	CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error)
	CreateRoute(routeRequest cfclient.RouteRequest) (route cfclient.Route, err error)
	DeleteApp(guid string) error
	GetAppBits(appGUID string) (r io.ReadCloser, err error)
	GetAppByGuid(guid string) (app cfclient.App, err error)
	GetAppInstances(guid string) (instances map[string]cfclient.AppInstance, err error)
	GetAppRoutes(appGUID string) ([]cfclient.Route, error)
	GetInfo() (*cfclient.Info, error)
	GetOrgByName(name string) (cfclient.Org, error)
	GetSharedDomainByName(name string) (cfclient.SharedDomain, error)
	GetSpaceByName(spaceName string, orgGUID string) (space cfclient.Space, err error)
	ListApps() ([]cfclient.App, error)
	ListBuildpacks() ([]cfclient.Buildpack, error)
	ListDomains() ([]cfclient.Domain, error)
	ListOrgQuotas() ([]cfclient.OrgQuota, error)
	ListOrgs() ([]cfclient.Org, error)
	ListRoutes() ([]cfclient.Route, error)
	ListSecGroups() ([]cfclient.SecGroup, error)
	ListServiceBindings() ([]cfclient.ServiceBinding, error)
	ListServiceBrokers() ([]cfclient.ServiceBroker, error)
	ListServiceInstances() ([]cfclient.ServiceInstance, error)
	ListServiceKeys() ([]cfclient.ServiceKey, error)
	ListServicePlanVisibilities() ([]cfclient.ServicePlanVisibility, error)
	ListServicePlans() ([]cfclient.ServicePlan, error)
	ListServices() ([]cfclient.Service, error)
	ListSharedDomains() ([]cfclient.SharedDomain, error)
	ListSpaces() ([]cfclient.Space, error)
	ListSpaceQuotas() ([]cfclient.SpaceQuota, error)
	ListStacks() ([]cfclient.Stack, error)
	ListUserProvidedServiceInstances() ([]cfclient.UserProvidedServiceInstance, error)
	ListUsers() (cfclient.Users, error)
	UpdateApp(appGUID string, aur cfclient.AppUpdateResource) (cfclient.UpdateResponse, error)
	UploadAppBits(file io.Reader, appGUID string) error
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)
	ListOrgsByQuery(query url.Values) ([]cfclient.Org, error)
}

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
