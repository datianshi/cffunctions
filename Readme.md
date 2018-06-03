# Functions wrapper

with [cf-client](https://github.com/cloudfoundry-community/go-cfclient)

Want **Less Imperative Style**

## Examples

## Print each org name

```
var printOrgName functions.OrgAction
printOrgName = func(org functions.Org) error {
  fmt.Println(org.ORG.Name)
  return nil
}
cf.EachOrg().Action(printOrgName)()
```

## Print each org name in parallell

```
cf.EachOrg().Parallel(printOrgName)()
```

## Print only system space name

```
var onlySystem functions.SpaceFilter
onlySystem = func(space functions.Space) bool {
  if space.SPACE.Name == "system" {
    return true
  }
  return false
}
_, err = cf.EachSpace().Filter(onlySystem).Action(printSpaceName)()
```

## Advance usage: Copy each org to another one

```
var createOrg functions.OrgAction = func(org functions.Org) error {
  _, err = client2.GetOrgByName(org.ORG.Name)
  if err != nil {
    _, err = client2.CreateOrg(cfclient.OrgRequest{Name: org.ORG.Name})
  } else {
    fmt.Printf("Org already exists: %s", org.ORG.Name)
    fmt.Println()
  }
  return err
}
cf.EachOrg().Action(createOrg)()
```

Refer [main.go](driver/main.go) for more examples
