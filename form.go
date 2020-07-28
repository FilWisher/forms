package main

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

func HTMLType(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Bool:
		return "checkbox", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "number", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number", nil
	case reflect.Ptr:
		return HTMLType(t.Elem())
	case reflect.String:
		return "text", nil
	default:
		return "", fmt.Errorf("HTMLType: no corresponding type for %v", t)
	}
}

func isValidStructPointer(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Ptr && v.Elem().IsValid() && v.Elem().Type().Kind() == reflect.Struct
}

func parseOpt(opts []string, opt string) string {
	for _, o := range opts {
		if strings.HasPrefix(o, opt) {
			return strings.TrimPrefix(o, opt)
		}
	}
	return ""
}

func Render(d interface{}) (template.HTML, error) {
	parts, err := RenderEach(d)
	if err != nil {
		return "", err
	}

	var out template.HTML
	for _, part := range parts {
		out += part + "\n"
	}

	return out, nil
}

func RenderEach(d interface{}) ([]template.HTML, error) {
	return render(reflect.ValueOf(d))
}

func render(v reflect.Value) ([]template.HTML, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("schema: interface must be a struct")
	}
	t := v.Type()

	var out []template.HTML

	for i := 0; i < v.NumField(); i++ {
		name, opts := fieldAlias(t.Field(i), "form")
		if name == "-" {
			continue
		}

		// Encode struct pointer types if the field is a valid pointer and a struct.
		if isValidStructPointer(v.Field(i)) {
			vs, err := render(v.Field(i).Elem())
			if err != nil {
				return nil, err
			}

			out = append(out, vs...)
			continue
		}

		if v.Field(i).Type().Kind() == reflect.Struct {
			vs, err := render(v.Field(i))
			if err != nil {
				return nil, err
			}
			out = append(out, vs...)
			continue
		}

		if v.Field(i).Type().Kind() == reflect.Slice {
			return nil, errors.New("form.Render: cannot render slice types")
		}

		typ := parseOpt(opts, "type=")

		// No type provided, look up default
		if typ == "" {
			def, err := HTMLType(v.Field(i).Type())
			if err != nil {
				return nil, err
			}
			typ = def
		}

		input := fmt.Sprintf("<input type=\"%s\" name=\"%s\">", typ, name)
		out = append(out, template.HTML(input))
	}

	return out, nil
}

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, []string) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// fieldAlias parses a field tag to get a field alias.
func fieldAlias(field reflect.StructField, tagName string) (alias string, options []string) {
	if tag := field.Tag.Get(tagName); tag != "" {
		alias, options = parseTag(tag)
	}
	if alias == "" {
		alias = field.Name
	}
	return alias, options
}
