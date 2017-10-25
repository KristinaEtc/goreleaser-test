package goreleaser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goreleaser/archive"
	zglob "github.com/mattn/go-zglob"
	"golang.org/x/sync/errgroup"
)

// Archive archives buildable files
func Archive(ctx *Context) error {
	var g errgroup.Group
	for platform, binaries := range ctx.Binaries {
		log.Printf("binaries=%+v\n", binaries)
		platform := platform
		binaries := binaries
		g.Go(func() error {
			return create(ctx, platform, binaries)
		})
	}
	return g.Wait()
}

func create(ctx *Context, platform string, groups map[string][]Binary) error {
	fmt.Printf("debug [create] groups=%+v", groups)
	for folder, binaries := range groups {
		log.Println("folder=", folder)

		var format = archiveformatFor(ctx, platform)

		archivePath := filepath.Join(dirWithGitHubReleases, folder+"."+format)
		archiveFile, err := os.Create(archivePath)
		fmt.Println("folder=", folder, binaries, format, archivePath, archiveFile)

		if err != nil {
			return fmt.Errorf("failed to create directory %s: %s", archivePath, err.Error())
		}
		defer func() {
			if e := archiveFile.Close(); e != nil {
				log.Println("archive", archivePath, "failed to close file: ", e.Error())
			}
		}()
		log.Println("archive", archivePath, "creating")

		var a = archive.New(archiveFile)
		defer func() {
			if e := a.Close(); e != nil {
				log.Println("archive", archivePath, "failed to close archive: ", e.Error())
			}
		}()

		files, err := findFiles(ctx)
		if err != nil {
			return fmt.Errorf("failed to find files to archive: %s", err.Error())
		}
		for _, f := range files {
			if err = a.Add(wrap(ctx, f, folder), f); err != nil {
				return fmt.Errorf("failed to add %s to the archive: %s", f, err.Error())
			}
		}
		for _, binary := range binaries {
			if err := a.Add(wrap(ctx, binary.Name, folder), binary.Path); err != nil {
				return fmt.Errorf("failed to add %s -> %s to the archive: %s", binary.Path, binary.Name, err.Error())
			}
		}
		ctx.AddArtifact(archivePath)
	}
	return nil
}

// archiveformatFor return the archive format, considering overrides and all that
func archiveformatFor(ctx *Context, platform string) string {
	for _, override := range ctx.Config.Archive.FormatOverrides {
		fmt.Println("\n\n[DEBUG] override", override, "platform=", platform)
		if strings.HasPrefix(platform, override.Goos) {
			return override.Format
		}
	}
	return ctx.Config.Archive.Format
}

func findFiles(ctx *Context) (result []string, err error) {
	for _, glob := range ctx.Config.Archive.Files {
		files, err := zglob.Glob(glob)
		if err != nil {
			return result, fmt.Errorf("globbing failed for pattern %s: %s", glob, err.Error())
		}
		result = append(result, files...)
	}
	return
}

// Wrap archive files with folder if set in config.
func wrap(ctx *Context, name, folder string) string {
	if ctx.Config.Archive.WrapInDirectory {
		return filepath.Join(folder, name)
	}
	return name
}
