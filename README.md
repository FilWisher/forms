# forms

Easily generate HTML form inputs from your types. Parse HTTP request parameters
back into types (like https://github.com/gorilla/schema).

```
inputs, _ := form.Render(NewUserForm{})
tmpl.Execute(os.Stdout, map[string]interface{}{
    "form": inputs,
})
```

The main function is `forms.Render(interface{}) (template.HTML, error)`. This
returns a set of HTML inputs that can be nested directly in your
`html/template.Template`.

The input types have sensible defaults but can be overwridden with struct tags
or with explicit `Options`:

```
form.RenderOpts(NewUserForm{}, map[string]Options{
    "Username": form.Options{Name: "username", ID: "username-field"},
})
```

Nested structs are supported and field names are namespaced. These types:

```
type Address struct {
    Line1 string
    Line2 string
    City string
}

type Person struct {
    Name string
    Address Address
}
```

will generate this form:

```
<input type='text' name='Name'>
<input type='text' name='Address.Line1'>
<input type='text' name='Address.Line2'>
<input type='text' name='Address.City'>
```

Slices are supported. These types:

```
type Pet struct { Name string }
type Person struct {
    Pets []Pet
}

person := &Person{
    Pets: []Pet{{Name: "hector"},{Name: "biggles"}},
}
```

will generate this form:

```
<input type='text' name='Person.Pets.1.Name' value='hector'>
<input type='text' name='Person.Pets.2.Name' value='biggles'>
```

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

