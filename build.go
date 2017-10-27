package goreleaser

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/KristinaEtc/goreleaser/buildtarget"
	"github.com/goreleaser/goreleaser/config"
)

var binariesLock sync.Mutex

type ldflagsData struct {
	Date    string
	Tag     string
	Commit  string
	Version string
}

func RunBuild(conf *config.Project) error {
	var ctx = &Context{
		Config: *conf,
	}
	setDefaultValues(ctx)
	for _, build := range ctx.Config.Builds {
		log.Println("INFO [RunBuld] building")
		if err := runBuild(ctx, &build); err != nil {
			return err
		}
	}
	return nil
}

func setDefaultValues(ctx *Context) {
	log.Println("DEBUG set default values")
}

func runBuild(ctx *Context, build *config.Build) error {
	log.Println("[DEBUG] runBuild")
	log.Println("len=", len(buildtarget.All(*build)))
	for _, target := range buildtarget.All(*build) {
		fmt.Printf("[DEBUG] target=%+v\n", target)
		target := target
		build := build

		doBuild(ctx, *build, target)
	}
	return nil
}

func doBuild(ctx *Context, build config.Build, target buildtarget.Target) error {
	folder, err := ForName(ctx, target)
	if err != nil {
		return err
	}
	var binaryName = build.Binary + buildtarget.For(target)
	var prettyName = binaryName
	if ctx.Config.Archive.Format == "binary" {
		binaryName, err = ForBuild(ctx, build, target)
		if err != nil {
			return err
		}
		log.Println("binaryName=", binaryName)
		//binaryName = binaryName + extentionFor(target)
		//log.Println("binaryName2=", binaryName)

	}
	var binary = filepath.Join(ctx.Config.Dist, folder, binaryName)
	log.Println("binary=", binary)
	ctx.AddBinary(target.String(), folder, prettyName, binary)
	cmd := []string{"go", "build"}
	if build.Flags != "" {
		cmd = append(cmd, strings.Fields(build.Flags)...)
	}
	flags, err := ldflags(ctx, build)
	if err != nil {
		return err
	}
	cmd = append(cmd, "-ldflags="+flags, "-o", binary, build.Main)
	return run(target, cmd, build.Env)
}

func ldflags(ctx *Context, build config.Build) (string, error) {
	var data = ldflagsData{
		Commit:  ctx.Git.Commit,
		Tag:     ctx.Git.CurrentTag,
		Version: ctx.Version,
		Date:    time.Now().UTC().Format(time.RFC3339),
	}
	var out bytes.Buffer
	t, err := template.New("ldflags").Parse(build.Ldflags)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, data)
	return out.String(), err
}

// AddBinary adds a built binary to the current context
func (ctx *Context) AddBinary(platform, folder, name, path string) {
	binariesLock.Lock()
	defer binariesLock.Unlock()
	if ctx.Binaries == nil {
		ctx.Binaries = map[string]map[string][]Binary{}
	}
	if ctx.Binaries[platform] == nil {
		ctx.Binaries[platform] = map[string][]Binary{}
	}
	ctx.Binaries[platform][folder] = append(
		ctx.Binaries[platform][folder],
		Binary{
			Name: name,
			Path: path,
		},
	)
}

func run(target buildtarget.Target, command, env []string) error {
	var cmd = exec.Command(command[0], command[1:]...)
	env = append(env, target.Env()...)

	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)
	log.Println("INFO [run] running")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println("ERR [run] failed: ", err.Error())
		return fmt.Errorf("build failed for %s:\n%v", target.String(), string(out))
	}
	return nil
}
