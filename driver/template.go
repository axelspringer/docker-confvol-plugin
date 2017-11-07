package driver

import (
	"bytes"
	"text/template"
)

// ConfTemplate data struct
type ConfTemplate struct {
	store      Store
	funcHelper template.FuncMap
}

func (ct *ConfTemplate) Parse(t string, opt interface{}) (string, error) {
	tmpl, err := template.New("conf_template").Funcs(ct.funcHelper).Parse(t)

	if err != nil {
		return "", err
	}

	var execBuffer bytes.Buffer
	if e := tmpl.Execute(&execBuffer, opt); e != nil {
		return "", e
	}

	return execBuffer.String(), nil
}

func (ct *ConfTemplate) ParseFromStore(k string) (string, error) {
	tmpData, err := ct.store.Get(k)

	if err != nil {
		return "", err
	}

	return ct.Parse(string(tmpData.Value), nil)
}

// NewTemplate create a new configuration template
func NewTemplate(s Store) *ConfTemplate {
	t := &ConfTemplate{
		store:      s,
		funcHelper: template.FuncMap{},
	}

	t.funcHelper["StoreGet"] = func(k string) string {
		if t.store == nil {
			return ""
		}

		entry, err := t.store.Get(k)
		if err != nil {
			return ""
		}

		return string(entry.Value)
	}

	t.funcHelper["StoreList"] = func(k string) []string {
		ret := []string{}
		if t.store == nil {
			return ret
		}

		entryList, _ := t.store.List(k)
		for _, kv := range entryList {
			ret = append(ret, string(kv.Value))
		}

		return ret
	}

	return t
}
