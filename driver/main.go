package main

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	functions "github.com/datianshi/cffunctions"
)

func main() {
	config := &cfclient.Config{
		ApiAddress:        "https://api.pas.nsx-t.shaozhenpcf.com",
		Username:          "admin",
		Password:          "hc9TO6_KThKb3WVhrX68B1uIusyLLvQM",
		SkipSslValidation: true,
	}
	client, err := cfclient.NewClient(config)
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to connect to cloudfoundry: %s", err.Error()))
		return
	}
	cf := &functions.CloudFoundry{
		Api: client,
	}

	//Print Each Org Name
	fmt.Println("Print Each Org Names:")
	var printOrgName functions.OrgAction
	printOrgName = func(org functions.Org) error {
		fmt.Println(org.ORG.Name)
		return nil
	}
	_, err = cf.EachOrg().Action(printOrgName)()
	if err != nil {
		fmt.Println(err)
	}

	//Print Each Org Name in parallel
	fmt.Println("Print Each Org Names in parallel:")
	_, err = cf.EachOrg().Parallel(printOrgName)()
	if err != nil {
		fmt.Println(err)
	}

	//Print Each Space Name
	fmt.Println("Print Each Space Names:")
	var printSpaceName functions.SpaceAction
	printSpaceName = func(space functions.Space) error {
		fmt.Println(space.SPACE.Name)
		return nil
	}
	_, err = cf.EachSpace().Action(printSpaceName)()
	if err != nil {
		fmt.Println(err)
	}

	//Filter Space based on Name - Only Print if Space name is 'system'
	fmt.Println("Only Print System Space Name:")
	var onlySystem functions.SpaceFilter
	onlySystem = func(space functions.Space) bool {
		if space.SPACE.Name == "system" {
			return true
		}
		return false
	}
	_, err = cf.EachSpace().Filter(onlySystem).Action(printSpaceName)()
	if err != nil {
		fmt.Println(err)
	}

	//Advanced Usage: Copy Orgs to another Cloudfoundry Foundation

	fmt.Println("copy to orgs to another foundation")
	config2 := &cfclient.Config{
		ApiAddress:        "https://api.pks.nsx.shaozhenpcf.com",
		Username:          "admin",
		Password:          "Vu056X1ePqy1q-v19rfTUS9YJLkdJeWd",
		SkipSslValidation: true,
	}
	client2, err := cfclient.NewClient(config2)
	if err != nil {
		fmt.Println(fmt.Errorf("Failed to connect to cloudfoundry: %s", err.Error()))
		return
	}

	var createOrg functions.OrgAction = func(org functions.Org) error {
		_, err = client2.GetOrgByName(org.ORG.Name)
		if err != nil {
			_, err = client2.CreateOrg(cfclient.OrgRequest{Name: org.ORG.Name})
			fmt.Println(err)
		} else {
			err = fmt.Errorf("Org already exists: %s", org.ORG.Name)
		}
		return err
	}

	_, err = cf.EachOrg().Parallel(createOrg)()
	if err != nil {
		fmt.Println(err)
	}

}
