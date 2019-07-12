// update-examples parses the README.md files in the examples directory
// and updates the code snippets with minimal versions of the shaders
// since there is not a good way to embed files into README.md files
// on GitHub.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	h2RE = regexp.MustCompile(`\n##\s+`)
)

type irmfFile struct {
	name     string
	contents string
}

func main() {
	readmeByPath := map[string]string{}
	irmfByPath := map[string]map[string]string{}
	if err := filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("path=%q, err=%v", path, err)
		}
		if info.IsDir() {
			irmfByPath[path] = map[string]string{}
			return nil
		}
		if info.Name() == "README.md" {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("ReadFile(%q): %v", path, err)
			}
			readmeByPath[filepath.Dir(path)] = string(buf)
			return nil
		}
		if strings.HasSuffix(path, ".irmf") {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("ReadFile(%q): %v", path, err)
			}
			dir := filepath.Dir(path)
			base := filepath.Base(path)
			irmfByPath[dir][base] = removeExtraFields(string(buf))
		}
		return nil
	}); err != nil {
		log.Fatalf("filepath.Walk: %v", err)
	}

	for k, v := range readmeByPath {
		processReadme(k, v, irmfByPath[k])
	}
}

func removeExtraFields(s string) string {
	lines := strings.Split(s, "\n")
	inJSON := true
	var out []string
	for i, line := range lines {
		if i == 0 {
			out = append(out, line)
			continue
		}
		if !inJSON {
			out = append(out, line)
			continue
		}
		if line == "}*/" {
			inJSON = false
			out = append(out, line)
			continue
		}
		for _, key := range fieldsToKeep {
			tag := `"` + key + `": `
			index := strings.Index(line, tag)
			if index >= 0 {
				out = append(out, "  "+key+": "+line[len(tag)+index:])
				break
			}
		}
	}
	return strings.Join(out, "\n")
}

func processReadme(path, buf string, irmfs map[string]string) {
	log.Printf("Processing %v/README.md...", path)
	log.Printf("Found %v .irmf files...", len(irmfs))

	parts := h2RE.Split(buf, -1)
	log.Printf("Found %v ## sections...", len(parts))
	for i, v := range parts {
		if i == 0 {
			continue
		}
		index := strings.Index(v, ".irmf")
		filename := v[:index+5]
		glsl, ok := irmfs[filename]
		if !ok {
			log.Fatalf("Could not find file %v", filename)
		}
		glslIndex := strings.Index(v, "```glsl")
		if glslIndex < 0 {
			log.Fatalf("Unable to find ```glsl...``` in %v", v)
		}
		parts[i] = "## " + v[0:glslIndex+8] + glsl + "```\n\n" + tryMessage(path, filename)
	}
	parts = append(parts, licenseText)

	outbuf := []byte(strings.Join(parts, "\n"))
	if err := ioutil.WriteFile(path+"/README.md", outbuf, 0644); err != nil {
		log.Fatalf("WriteFile: %v", err)
	}
}

func tryMessage(path, filename string) string {
	return fmt.Sprintf(`* Try loading [%v](https://gmlewis.github.io/irmf-editor/?s=github.com/gmlewis/irmf/blob/master/%v/%v) now in the experimental IRMF editor!`+"\n", filename, path, filename)
}

var fieldsToKeep = []string{
	"irmf",
	"materials",
	"max",
	"min",
	"units",
}

var licenseText = `----------------------------------------------------------------------

# License

Copyright 2019 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
