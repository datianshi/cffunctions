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
	var actionDriver OrgAction
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
			actionDriver = func(org Org) error {
				catchOrgs = append(catchOrgs, org.ORG.Name)
				return nil
			}
		})
		Context("Sequential Execution", func() {
			BeforeEach(func() {
				_, err = cf.EachOrg().Action(actionDriver)()
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
		})
		Context("Parallel Execution with org2 and org3 failed", func() {
			BeforeEach(func() {
				actionDriver = func(org Org) error {
					if org.ORG.Name == "org3" {
						return errors.New("org3 failed!")
					}
					if org.ORG.Name == "org2" {
						return errors.New("org2 failed!")
					}
					catchOrgs = append(catchOrgs, org.ORG.Name)
					return nil
				}
				_, err = cf.EachOrg().Parallel(actionDriver)()
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
