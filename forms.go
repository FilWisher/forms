package forms

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

type Options struct {
	Type  string
	Class string
	ID    string
	Name  string
}

func parseOpts(opts []string) (Options, error) {
	var options Options
	for _, opt := range opts {
		parts := strings.Split(opt, "=")
		if len(parts) != 2 {
			return options, fmt.Errorf("invalid option format: %v", opt)
		}
		switch parts[0] {
		case "type":
			options.Type = parts[1]
		case "class":
			options.Class = parts[1]
		case "id":
			options.ID = parts[1]
		default:
			return options, fmt.Errorf("invalid option: %v", parts[0])
		}
	}
	return options, nil
}

func parseOpt(opts []string, opt string) string {
	for _, o := range opts {
		if strings.HasPrefix(o, opt) {
			return strings.TrimPrefix(o, opt)
		}
	}
	return ""
}

func RenderOpts(d interface{}, opts map[string]Options) (template.HTML, error) {
	parts, err := RenderEachOpts(d, opts)
	if err != nil {
		return "", err
	}

	var out template.HTML
	for _, part := range parts {
		out += part + "\n"
	}

	return out, nil
}

func Render(d interface{}) (template.HTML, error) {
	return RenderOpts(d, nil)
}

func RenderEachOpts(d interface{}, opts map[string]Options) ([]template.HTML, error) {
	return render(reflect.ValueOf(d), opts)
}

func RenderEach(d interface{}) ([]template.HTML, error) {
	return RenderEachOpts(d, nil)
}

func renderInput(opts Options) template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<input type=\"%s\" name=\"%s\"", opts.Type, opts.Name))
	if opts.Class != "" {
		sb.WriteString(fmt.Sprintf(" class=\"%s\"", opts.Class))
	}
	if opts.ID != "" {
		sb.WriteString(fmt.Sprintf(" id=\"%s\"", opts.ID))
	}
	sb.WriteString(">")
	return template.HTML(sb.String())
}

func mergeOptions(opts Options, defaults Options) Options {
	if opts.Name == "" {
		opts.Name = defaults.Name
	}
	if opts.ID == "" {
		opts.ID = defaults.ID
	}
	if opts.Class == "" {
		opts.Class = defaults.Class
	}
	if opts.Type == "" {
		opts.Type = defaults.Type
	}
	return opts
}

func render(v reflect.Value, optsMap map[string]Options) ([]template.HTML, error) {
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
			vs, err := render(v.Field(i).Elem(), optsMap)
			if err != nil {
				return nil, err
			}

			out = append(out, vs...)
			continue
		}

		if v.Field(i).Type().Kind() == reflect.Struct {
			vs, err := render(v.Field(i), optsMap)
			if err != nil {
				return nil, err
			}
			out = append(out, vs...)
			continue
		}

		if v.Field(i).Type().Kind() == reflect.Slice {
			return nil, errors.New("form.Render: cannot render slice types")
		}

		// Prefer supplied options. Fallback to using tag-based options.
		var options Options
		if optsMap != nil {
			var ok bool
			options, ok = optsMap[name]
			if !ok {
				options = Options{Name: name}
			}
		}

		tagOptions, err := parseOpts(opts)
		if err != nil {
			return nil, fmt.Errorf("Render: %w", err)
		}

		options = mergeOptions(options, tagOptions)

		// No type provided, look up default
		if options.Type == "" {
			def, err := HTMLType(v.Field(i).Type())
			if err != nil {
				return nil, err
			}
			options.Type = def
		}

		out = append(out, renderInput(options))
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
