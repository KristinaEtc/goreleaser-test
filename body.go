package goreleaser

import (
	"bytes"
	"os/exec"
	"text/template"
)

const bodyTemplate = `{{ .ReleaseNotes }}

---
Automated with [GoReleaser](https://github.com/goreleaser)
Built with {{ .GoVersion }}`

func describeBody(ctx *Context) (bytes.Buffer, error) {
	bts, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return bytes.Buffer{}, err
	}
	return describeBodyVersion(ctx, string(bts))
}

func describeBodyVersion(ctx *Context, version string) (bytes.Buffer, error) {
	var out bytes.Buffer
	var template = template.Must(template.New("release").Parse(bodyTemplate))
	err := template.Execute(&out, struct {
		ReleaseNotes, GoVersion string
		//	DockerImages            []string
	}{
		ReleaseNotes: ctx.ReleaseNotes,
		GoVersion:    version,
		//DockerImages: ctx.Dockers,
	})
	return out, err
}
