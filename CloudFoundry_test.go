package cffunctions_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/datianshi/cffunctions/cffunctionsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/datianshi/cffunctions"
)

var _ = Describe("CloudFoundry", func() {
	var api *cffunctionsfakes.FakeAPI
	var cf *CloudFoundry
	var err error
	var catchOrgs []string
	var originOrgs []string
	var actionDriver Action
	BeforeEach(func() {
		api = &cffunctionsfakes.FakeAPI{}
		cf = &CloudFoundry{Api: api}
	})

	Context("Given 5 Orgs Returns", func() {
		BeforeEach(func() {
			originOrgs = []string{"org1", "org2", "org3", "org4", "org5"}
			api.ListOrgsStub = func() ([]cfclient.Org, error) {
				orgs := []cfclient.Org{
					cfclient.Org{Name: originOrgs[0]},
					cfclient.Org{Name: originOrgs[1]},
					cfclient.Org{Name: originOrgs[2]},
					cfclient.Org{Name: originOrgs[3]},
					cfclient.Org{Name: originOrgs[4]},
				}
				return orgs, nil
			}

			catchOrgs = make([]string, 0)
			actionDriver = func(org CFObject) error {
				catchOrgs = append(catchOrgs, org.GetName())
				return nil
			}
		})
		Context("Sequential Execution", func() {
			BeforeEach(func() {
				_, err = cf.Each(OrgFunc).Action(actionDriver)()
			})
			It("Should have no error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("Should Catch 5 orgs", func() {
				Ω(len(catchOrgs)).Should(Equal(5))
			})
			It("Should Equal to the originOrgs", func() {
				Ω(catchOrgs).Should(Equal(originOrgs))
			})

			Context("Given a filter to print only org1", func() {
				BeforeEach(func() {
					catchOrgs = make([]string, 0)
					var filter Filter = func(o CFObject) bool {
						if o.GetName() == "org1" {
							return true
						}
						return false
					}
					_, err = cf.Each(OrgFunc).Filter(filter).Action(actionDriver)()
				})
				It("Should have no error", func() {
					Ω(err).ShouldNot(HaveOccurred())
				})
				It("Should Catch only 1 org", func() {
					Ω(len(catchOrgs)).Should(Equal(1))
				})
			})
		})
		Context("Parallel Execution with org2 and org3 failed", func() {
			BeforeEach(func() {
				actionDriver = func(object CFObject) error {
					if object.GetName() == "org3" {
						return errors.New("org3 failed!")
					}
					if object.GetName() == "org2" {
						return errors.New("org2 failed!")
					}
					catchOrgs = append(catchOrgs, object.GetName())
					return nil
				}
				_, err = cf.Each(OrgFunc).Parallel(actionDriver)()
			})
			It("Should have error", func() {
				Ω(err).Should(HaveOccurred())
			})
			It("Should Catch 5 orgs", func() {
				Ω(len(catchOrgs)).Should(Equal(3))
			})
			It("Should Not Equal to the originOrgs", func() {
				Ω(catchOrgs).ShouldNot(Equal(originOrgs))
			})
		})
	})

})
