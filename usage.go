package main

import (
	"github.com/filwisher/forms"

	"html/template"
	"os"
)

type NewUserForm struct {
	Username string `form:"username,type=text"`
	Password string `form:"password,type=password"`
	Confirm  string `form:"password,type=confirm"`
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

func main() {

	form, err := Render(NewUserForm{})
	if err != nil {
		panic(err)
	}

	base.Execute(os.Stdout, map[string]interface{}{
		"form": form,
	})

}
