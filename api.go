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
