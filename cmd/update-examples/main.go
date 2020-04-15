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
	"sort"
	"strings"
)

var (
	h2RE = regexp.MustCompile(`\n##\s+`)
)

func main() {
	readmeByPath := map[string]string{}
	irmfByPath := map[string]map[string]string{}
	stlFileSizesByPath := map[string]map[string]int64{}
	dlpFileSizesByPath := map[string]map[string]int64{}
	if err := filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("path=%q, err=%v", path, err)
		}
		if info.IsDir() {
			irmfByPath[path] = map[string]string{}
			stlFileSizesByPath[path] = map[string]int64{}
			dlpFileSizesByPath[path] = map[string]int64{}
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
		if strings.HasSuffix(path, ".stl") {
			dir := filepath.Dir(path)
			base := filepath.Base(path)
			stlFileSizesByPath[dir][base] = info.Size()
			return nil
		}
		if strings.HasSuffix(path, ".cbddlp") {
			dir := filepath.Dir(path)
			base := filepath.Base(path)
			dlpFileSizesByPath[dir][base] = info.Size()
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
		processReadme(k, v, irmfByPath[k], stlFileSizesByPath[k], dlpFileSizesByPath[k])
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

func processReadme(path, buf string, irmfs map[string]string, stlFileSizes map[string]int64, dlpFileSizes map[string]int64) {
	log.Printf("Processing %v/README.md ...", path)
	log.Printf("Found %v .irmf files...", len(irmfs))
	log.Printf("Found %v .stl files...", len(stlFileSizes))
	log.Printf("Found %v .cbddlp files...", len(dlpFileSizes))

	licenseText := newLicenseText

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

		if j := strings.Index(v, "-----"); j >= 0 {
			licenseText = v[j:] // Preserve year of original license text.
		}

		parts[i] = "## " + v[0:glslIndex+8] + glsl + "```\n\n" + tryMessage(path, filename)

		if len(stlFileSizes) > 0 {
			parts[i] += addSTLs(filename, stlFileSizes)
		}

		if len(dlpFileSizes) > 0 {
			parts[i] += addDLPs(filename, dlpFileSizes)
		}
	}
	parts = append(parts, licenseText)

	outbuf := []byte(strings.Join(parts, "\n"))
	if err := ioutil.WriteFile(path+"/README.md", outbuf, 0644); err != nil {
		log.Fatalf("WriteFile: %v", err)
	}
}

func addSTLs(filename string, stlFileSizes map[string]int64) string {
	var lines []string

	// Strip off the ".irmf"
	filename = strings.TrimSuffix(filename, ".irmf")

	for k, v := range stlFileSizes {
		if strings.HasPrefix(k, filename+"-mat") {
			lines = append(lines, fmt.Sprintf("  - [%v](%v) (%v bytes)", k, k, v))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	header := "* Here is a crude STL approximation of this model\n  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer)"
	if len(lines) > 1 {
		sort.Strings(lines)
		header += "\n  (one STL file per material)"
	}

	return "\n" + header + ":\n" + strings.Join(lines, "\n") + "\n"
}

func addDLPs(filename string, dlpFileSizes map[string]int64) string {
	var lines []string

	// Strip off the ".irmf"
	filename = strings.TrimSuffix(filename, ".irmf")

	for k, v := range dlpFileSizes {
		if strings.HasPrefix(k, filename+"-mat") {
			lines = append(lines, fmt.Sprintf("  - [%v](%v) (%v bytes)", k, k, v))
		}
	}

	if len(lines) == 0 {
		return ""
	}

	header := "* Here is a voxel approximation of this model\n  using [irmf-slicer](https://github.com/gmlewis/irmf-slicer)"
	if len(lines) > 1 {
		sort.Strings(lines)
		header += "\n  (one .cbddlp file per material)"
	}

	return "\n" + header + ":\n" + strings.Join(lines, "\n") + "\n"
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

var newLicenseText = `----------------------------------------------------------------------

# License

Copyright 2020 Glenn M. Lewis. All Rights Reserved.

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
