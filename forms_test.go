package forms

import (
	"html/template"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

type NewUserForm struct {
	Username string `form:"username,type=text"`
	Password string `form:"password,type=password"`
	Confirm  string `form:"confirm,type=password"`
}

var base = template.Must(template.New("base").Parse(`
<!DOCTYPE html>
<html>
<head></head>
<body>
<form>
{{ .form }}
</form>
</body></html>`))

func TestRender(t *testing.T) {

	form, err := Render(NewUserForm{
		Username: "filwisher",
	})
	if err != nil {
		t.Errorf("error rendering form: %v", err)
	}

	_, err = html.ParseFragment(strings.NewReader(string(form)), nil)
	if err != nil {
		t.Errorf("invalid html: %v", err)
	}

	err = base.Execute(os.Stdout, map[string]interface{}{
		"form": form,
	})
	if err != nil {
		t.Errorf("error executing template: %v", err)
	}
}

func TestRenderOpts(t *testing.T) {

	form, err := RenderOpts(NewUserForm{
		Password: "hunter2",
	}, map[string]Options{
		"Username": Options{Name: "email", Type: "text"},
		"confirm":  Options{Name: "confirm-password", Type: "password"},
		"password": Options{Value: ""},
	})
	if err != nil {
		t.Errorf("error rendering form: %v", err)
	}

	_, err = html.ParseFragment(strings.NewReader(string(form)), nil)
	if err != nil {
		t.Errorf("invalid html: %v", err)
	}

	err = base.Execute(os.Stdout, map[string]interface{}{
		"form": form,
	})
	if err != nil {
		t.Errorf("error executing template: %v", err)
	}
}
