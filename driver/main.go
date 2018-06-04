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
	var printName functions.Action
	printName = func(o functions.CFObject) error {
		fmt.Println(o.GetName())
		return nil
	}
	_, err = cf.Each(functions.OrgFunc).Action(printName)()
	if err != nil {
		fmt.Println(err)
	}

	//Print Each Org Name in parallel
	fmt.Println("Print Each Org Names in parallel:")
	_, err = cf.Each(functions.OrgFunc).Parallel(printName)()
	if err != nil {
		fmt.Println(err)
	}

	//Print Each Space Name
	fmt.Println("Print Each Space Names:")
	_, err = cf.Each(functions.SpaceFunc).Action(printName)()
	if err != nil {
		fmt.Println(err)
	}

	//Filter Space based on Name - Only Print if Space name is 'system'
	fmt.Println("Only Print System Space Name:")
	var onlySystem functions.Filter
	onlySystem = func(o functions.CFObject) bool {
		if o.GetName() == "system" {
			return true
		}
		return false
	}
	_, err = cf.Each(functions.SpaceFunc).Filter(onlySystem).Action(printName)()
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

	var createOrg functions.Action = func(o functions.CFObject) error {
		_, err = client2.GetOrgByName(o.GetName())
		if err != nil {
			_, err = client2.CreateOrg(cfclient.OrgRequest{Name: o.GetName()})
			fmt.Println(err)
		} else {
			err = fmt.Errorf("Org already exists: %s", o.GetName())
		}
		return err
	}

	_, err = cf.Each(functions.OrgFunc).Parallel(createOrg)()
	if err != nil {
		fmt.Println(err)
	}

}
