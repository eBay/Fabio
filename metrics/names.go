package metrics

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Service struct {
	Service   string
	Host      string
	Path      string
	TargetURL *url.URL
}

// DefaultNames contains the default template for route metric names for backends that don't
// support tags.
const DefaultNames = "{{clean .Service}}.{{clean .Host}}.{{clean .Path}}.{{clean .TargetURL.Host}}"

// DefaulPrefix contains the default template for metrics prefix.
const DefaultPrefix = "{{clean .Hostname}}.{{clean .Exec}}"

// names stores the template for the route metric names.
var names *template.Template

func init() {
	// make sure names is initialized to something
	var err error
	if names, err = parseNames(DefaultNames); err != nil {
		panic(err)
	}
}

func (s Service) String() string {
	return s.Service
}

const DotSeparator = "."
const PipeSeparator = "|"
const RoutePrefix = "route"

func Flatten(name string, labels []string, separator string) string {
	if len(labels) == 0 {
		return name
	}
	return name + separator + strings.Join(labels, separator)
}

func Labels(labels, values []string, stringsprefix, fieldsep, recsep string) string {
	if len(labels) == 0 {
		return ""
	}
	var b strings.Builder
	_, _ = b.WriteString(stringsprefix)
	for i := range labels {
		if i > 0 {
			_, _ = b.WriteString(recsep)
		}
		_, _ = b.WriteString(labels[i])
		_, _ = b.WriteString(fieldsep)
		if i < len(values) {
			_, _ = b.WriteString(values[i])
		}
	}
	return b.String()
}

// parseNames parses the route metric name template.
func parseNames(tmpl string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"clean":      clean,
		"clean_prom": clean_prom,
	}
	t, err := template.New("names").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return nil, err
	}
	testURL, err := url.Parse("http://127.0.0.1:12345/")
	if err != nil {
		return nil, err
	}
	if _, err := TargetName("testservice", "test.example.com", "/test", testURL.String()); err != nil {
		return nil, err
	}
	return t, nil
}

// parsePrefix parses the prefix metric template
func parsePrefix(tmpl string) (string, error) {
	// Backward compatibility condition for old metrics.prefix parameter 'default'
	if tmpl == "default" {
		tmpl = DefaultPrefix
	}
	funcMap := template.FuncMap{
		"clean": clean,
	}
	t, err := template.New("prefix").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}
	host, err := hostname()
	if err != nil {
		return "", err
	}
	exe := filepath.Base(os.Args[0])

	b := new(bytes.Buffer)
	data := struct{ Hostname, Exec string }{host, exe}
	if err := t.Execute(b, &data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// clean creates safe names for graphite reporting by replacing
// some characters with underscores.
// TODO(fs): This may need updating for other metrics backends.
func clean(s string) string {
	if s == "" {
		return "_"
	}
	s = strings.Replace(s, ".", "_", -1)
	s = strings.Replace(s, ":", "_", -1)
	return strings.ToLower(s)
}

// clean_prom creates safe names for prometheus reporting by replacing
// some characters with underscores.
func clean_prom(s string) string {
	if s == "" {
		return "_"
	}
	s = strings.Replace(s, ".", "_", -1)
	s = strings.Replace(s, ":", "_", -1)
	s = strings.Replace(s, "-", "_", -1)
	return strings.ToLower(s)
}

// stubbed out for testing
var hostname = os.Hostname

// TargetName returns the metrics name from the given parameters.
func TargetName(service, host, path, target string) (string, error) {
	if names == nil {
		return "", nil
	}

	targetURL, err := url.Parse(target)

	if err != nil {
		return "", fmt.Errorf("error parsing URL %s: %w", target, err)
	}

	var name bytes.Buffer

	data := struct {
		Service, Host, Path string
		TargetURL           *url.URL
	}{service, host, path, targetURL}

	if err := names.Execute(&name, data); err != nil {
		return "", err
	}

	return name.String(), nil
}

// TargetNameWith - this is used for flat metrics backends (no tags support)
// in With() methods on target metrics.
func TargetNameWith(name string, values []string) (string, error) {

	if len(values)%2 == 1 {
		values = append(values, "unknown")
	}
	m := make(map[string]string)
	for i := 0; i < len(values); i += 2 {
		m[values[i]] = values[i+1]
	}

	n, err := TargetName(m["service"], m["host"], m["path"], m["target"])
	if err != nil {
		return "", err
	}
	// Handle .tx and .rx
	if i := strings.LastIndex(name, "."); i != -1 {
		n += name[i:]
	}
	return n, nil
}

func isRouteMetric(name string) bool { return strings.HasPrefix(name, RoutePrefix) }
