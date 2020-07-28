# forms

Easily generate HTML form inputs from your types.

```
inputs, _ := form.Render(NewUserForm{})
tmpl.Execute(os.Stdout, map[string]interface{}{
    "form": inputs,
})
```

The main function is `forms.Render(interface{}) (template.HTML, error)`. This
returns a set of HTML inputs that can be nested directly in your
`html/template.Template`.

The input types have sensible defaults but can be overwridden with struct tags.

## Usage

Full example

```
package main

import (
	"html/template"
	"os"

	"github.com/filwisher/forms"
)

type NewUser struct {
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
    </body>
</html>`))

func main() {

	form, err := forms.Render(NewUser{})
	if err != nil {
		panic(err)
	}

	base.Execute(os.Stdout, map[string]interface{}{
		"form": form,
	})
}
```
