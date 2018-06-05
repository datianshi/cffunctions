# Functions wrapper

with [cf-client](https://github.com/cloudfoundry-community/go-cfclient)

Want **Less Imperative Style**

## Examples

## Print each org name

```
var printName functions.Action
printName = func(o functions.CFObject) error {
  fmt.Println(o.GetName())
  return nil
}
_, err = cf.Each(functions.OrgFunc).Action(printName)()
if err != nil {
  fmt.Println(err)
}
```

## Print each org name in parallel

* Use GOROUTINE to serve the parallel processing
* Aggregate all the sub stream errors

```
cf.Each(functions.OrgFunc).Parallel(printName)()
```

## Print only system space name

```
var onlySystem functions.Filter
onlySystem = func(o functions.CFObject) bool {
  if o.GetName() == "system" {
    return true
  }
  return false
}
_, err = cf.Each(functions.SpaceFunc).Filter(onlySystem).Action(printName)()
```

## Advance usage: Copy each org to another one in parallel

```
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
```

Refer [main.go](driver/main.go) for more examples
