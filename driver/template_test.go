package driver_test

import (
	"errors"
	"regexp"
	"testing"

	. "github.com/axelspringer/docker-conf-volume/driver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTemplates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Template Suite")
}

type StoreMock struct {
	kvMap map[string]string
	get   func(key string) (*StoreKVPair, error)
	list  func(key string) ([]*StoreKVPair, error)
}

func (s *StoreMock) Get(p string) (*StoreKVPair, error) {
	if s.get != nil {
		return s.get(p)
	}

	if entry, ok := s.kvMap[p]; ok {
		return &StoreKVPair{Key: p, Value: []byte(entry)}, nil
	}

	return nil, errors.New("Key not found in store")
}

func (s *StoreMock) List(p string) ([]*StoreKVPair, error) {
	if s.list != nil {
		return s.list(p)
	}

	l := []*StoreKVPair{}
	validEntry := regexp.MustCompile("^" + p + "[^/]$")

	for k, v := range s.kvMap {
		if validEntry.MatchString(k) {
			l = append(l, &StoreKVPair{Key: k, Value: []byte(v)})
		}
	}

	return l, nil
}

func newStoreMock(kv *map[string]string) *StoreMock {
	var m *map[string]string = kv
	if m == nil {
		m = &map[string]string{}
	}

	s := &StoreMock{
		kvMap: *m,
	}

	return s
}

var _ = Describe("Template", func() {

	// Check the initialization state of the ConfTemplate
	Context("Initialization", func() {
		template := NewTemplate(nil)

		It("has the correct type", func() {
			Expect(template).NotTo(BeNil())
		})
	})

	// Check the template value completion
	Context("Fill template value", func() {
		template := NewTemplate(newStoreMock(nil))

		It("Template is initialized with a concrete store", func() {
			Expect(template).NotTo(BeNil())
		})

		It("Test for parse error", func() {
			output, err := template.Parse("{{ ", nil)
			Expect(err).Should(MatchError("template: conf_template:1: unexpected unclosed action in command"))
			Expect(output).Should(Equal(""))
		})

		It("Parse textonly string template", func() {
			output, err := template.Parse("☂textonly☂template☂", nil)
			Expect(err).To(BeNil())
			Expect(output).Should(Equal("☂textonly☂template☂"))
		})

		It("Use template context", func() {
			output, err := template.Parse("{{ .Context }}", struct{ Context string }{"qwerty"})
			Expect(err).To(BeNil())
			Expect(output).Should(Equal("qwerty"))
		})
	})

	// StoreGet template helper function
	Context("Test StoreGet helper", func() {

		It("should resolved to empty string when store is nil", func() {
			templateWithoutStore := NewTemplate(nil)

			output, err := templateWithoutStore.Parse("{{StoreGet \"/foo/bar/buzz\"}}", nil)
			Expect(err).To(BeNil())
			Expect(output).Should(Equal(""))
		})

		It("should resolved to empty string when the store can't deliver", func() {
			template := NewTemplate(newStoreMock(nil))

			output, err := template.Parse("{{StoreGet \"/foo/bar/buzz\"}}", nil)
			Expect(err).To(BeNil())
			Expect(output).Should(Equal(""))
		})

		It("should resolved to entry value on success", func() {
			sm := newStoreMock(nil)
			sm.kvMap["/foo/bar/buzz"] = "qwerty"

			template := NewTemplate(sm)

			output, err := template.Parse("{{StoreGet \"/foo/bar/buzz\"}}", nil)
			Expect(err).To(BeNil())
			Expect(output).Should(Equal("qwerty"))
		})
	})

	// StoreGet template helper function
	Context("Test StoreList helper", func() {

		It("should resolved to entry value on success", func() {
			sm := newStoreMock(nil)
			sm.kvMap["/foo/bar/A"] = "qwertyA"

			template := NewTemplate(sm)

			output, err := template.Parse("{{range StoreList \"/foo/bar/\"}}{{ . }}{{end}}", nil)
			Expect(err).To(BeNil())
			Expect(output).Should(Equal("qwertyA"))
		})

	})
})
