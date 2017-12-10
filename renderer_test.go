package webutil_test

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/noypi/router"
	"github.com/noypi/webutil"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestRendererCloneTemplates(t *testing.T) {
	assert := assertpkg.New(t)

	tpl1 := `base template: {{ .Value }}`
	tpl2 := `custom template: {{ .CustomValue }}, {{template "base" .}}`
	tpl3 := `layout: {{template "custom" .}}`

	t1 := template.Must(template.New("base").Parse(tpl1))
	t2 := template.Must(template.New("custom").Parse(tpl2))
	t3 := template.Must(template.New("layout").Parse(tpl3))

	type Base struct {
		TPLConfig string `webtpl:"name=base"`
		Value     string
	}
	type Custom struct {
		TPLConfig   string `webtpl:"name=custom"`
		CustomValue string
	}
	type Layout struct {
		TPLConfig   string `webtpl:"name=layout"`
		CustomValue string
	}

	tAll := template.New("")
	tAll = template.Must(tAll.AddParseTree(t1.Name(), t1.Tree))
	tAll = template.Must(tAll.AddParseTree(t2.Name(), t2.Tree))
	tAll = template.Must(tAll.AddParseTree(t3.Name(), t3.Tree))

	oRenderer := webutil.NewRenderer(&router.Context{}, tAll)

	base := Base{Value: "some base value"}
	custom := Custom{CustomValue: "some custom value"}
	layout := Layout{CustomValue: "some layout value"}

	// test 1
	buf := new(bytes.Buffer)
	data := map[string]interface{}{}
	webutil.MergePagesData(nil, data, base, custom)
	tpl, err := oRenderer.CloneTemplate(base, custom)
	assert.Nil(err)
	assert.Nil(tpl.ExecuteTemplate(buf, "custom", data))
	assert.Equal(`custom template: some custom value, base template: some base value`, buf.String())

	// test 2
	data = map[string]interface{}{}
	webutil.MergePagesData(nil, data, layout)
	tpl, err = oRenderer.CloneTemplate(layout)
	assert.Nil(err)
	buf = new(bytes.Buffer)
	err = tpl.ExecuteTemplate(buf, "layout", data)
	assert.NotNil(err)
	assert.Contains(err.Error(), `no such template "custom"`)

	// test 3
	data = map[string]interface{}{}
	webutil.MergePagesData(nil, data, layout, base, custom)
	tpl, err = oRenderer.CloneTemplate(layout, base, custom)
	assert.Nil(err)
	buf = new(bytes.Buffer)
	assert.Nil(tpl.ExecuteTemplate(buf, "layout", data))
	assert.Equal(`layout: custom template: some layout value, base template: some base value`, buf.String())
}
