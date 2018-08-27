package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/tools/cover"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// https://github.com/golang/go/blob/master/src/cmd/cover/html.go#L96
func percentCovered(p *cover.Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

// https://github.com/golang/go/blob/master/src/cmd/cover/func.go#L178
func findPkgs(profiles []*cover.Profile) (map[string]*Pkg, error) {
	pkgs := make(map[string]*Pkg)
	var list []string
	for _, profile := range profiles {
		if strings.HasPrefix(profile.FileName, ".") || filepath.IsAbs(profile.FileName) {
			continue
		}
		pkg := path.Dir(profile.FileName)
		if _, ok := pkgs[pkg]; !ok {
			pkgs[pkg] = nil
			list = append(list, pkg)
		}
	}

	goTool := filepath.Join(runtime.GOROOT(), "bin/go")
	cmd := exec.Command(goTool, append([]string{"list", "-e", "-json"}, list...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot run go list: %v\n%s", err, stderr.Bytes())
	}
	dec := json.NewDecoder(bytes.NewReader(stdout))
	for {
		var pkg Pkg
		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding go list json: %v", err)
		}
		pkgs[pkg.ImportPath] = &pkg
	}
	return pkgs, nil
}

// https://github.com/golang/go/blob/master/src/cmd/cover/func.go#L219
func findFile(pkgs map[string]*Pkg, file string) (string, error) {
	if strings.HasPrefix(file, ".") || filepath.IsAbs(file) {
		return file, nil
	}
	pkg := pkgs[path.Dir(file)]
	if pkg != nil {
		if pkg.Dir != "" {
			return filepath.Join(pkg.Dir, path.Base(file)), nil
		}
		if pkg.Error != nil {
			return "", errors.New(pkg.Error.Err)
		}
	}
	return "", fmt.Errorf("did not find package for %s in go list output", file)
}

// https://github.com/golang/go/blob/master/src/cmd/cover/func.go#L169
type Pkg struct {
	ImportPath string
	Dir        string
	Error      *struct {
		Err string
	}
}

// https://github.com/golang/go/blob/master/src/cmd/cover/html.go#L170
//
// modified to no 'set'
type templateData struct {
	Files []*templateFile
}

// https://github.com/golang/go/blob/master/src/cmd/cover/html.go#L175
//
// modified template.HTML to string
type templateFile struct {
	Name     string
	Body     string
	Coverage float64
}
