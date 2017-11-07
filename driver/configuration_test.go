package driver_test

import (
	"testing"

	. "github.com/axelspringer/docker-conf-volume/driver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}

var _ = Describe("Configuration", func() {

	// Check the initialization state of the Configuration
	Context("Initialization", func() {
		It("has the correct type", func() {
			Expect(NewConfiguration()).NotTo(BeNil())
		})
	})

	// Check error handling
	Context("Error handling", func() {
		It("can handle nonexisting file errors", func() {
			err := NewConfiguration().LoadFromFile("do/not/exist")
			Expect(err.Error()).Should(Equal("open do/not/exist: no such file or directory"))
		})

		It("can handle empty data errors", func() {
			err := NewConfiguration().LoadFromString("")
			Expect(err.Error()).Should(Equal("Loading empty json data"))
		})

		It("can handle syntactic json errors", func() {
			err := NewConfiguration().LoadFromString("{\"flag\":}")
			Expect(err.Error()).Should(Equal("invalid character '}' looking for beginning of value"))
		})
	})

	// Check error handling
	Context("Parsing json data", func() {
		validJSONConf := `{ 
			"driver": {
				"rootpath": "/tmp/confvol/"
			},
			"backend": {
				"type": "etcd"
			},
			"generator": {
				"disabled": true
			}
		}`

		It("can handle valid json data", func() {
			conf := NewConfiguration()
			err := conf.LoadFromString(validJSONConf)
			Expect(err).To(BeNil())

			Expect(conf.Driver.RootPath).To(Equal("/tmp/confvol/"))
			Expect(conf.Backend.Type).To(Equal("etcd"))
			Expect(conf.Generator.Disabled).To(Equal(true))
		})

		It("can verify configuration integrity", func() {
			conf := NewConfiguration()
			conf.Backend.Type = "consul"

			integer, errList := conf.CheckIntegrity()
			Expect(integer).To(Equal(false))
			Expect(len(errList)).To(Equal(3))

			Expect(errList[0].Error()).To(Equal("driver.rootpath directory did not exist"))
			Expect(errList[1].Error()).To(Equal("backend.type only supports 'etcd' at the moment"))
			Expect(errList[2].Error()).To(Equal("backend.endpoints is a neccessary field"))
		})

		It("can verify configuration integrity", func() {
			conf := NewConfiguration()
			conf.Backend.Endpoints = " 10.0.0.1,    10.0.0.2"
			Expect(conf.GetBackendEndpointList()).To(Equal([]string{"10.0.0.1", "10.0.0.2"}))
		})
	})

})
