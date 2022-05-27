package sln

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Sln struct {
	SolutionDir string
	ProjectList []Project
}

func NewSln(path string) (Sln, error) {
	var sln Sln
	var err error

	sln.SolutionDir, err = filepath.Abs(path)
	sln.SolutionDir = filepath.Dir(sln.SolutionDir)
	if err != nil {
		return sln, err
	}
	projectFiles, err := findAllProject(path)
	if err != nil {
		fmt.Println(err)
		return sln, err
	}
	if len(projectFiles) == 0 {
		return sln, errors.New("not found project file")
	}

	for _, path := range projectFiles {
		pro, err := NewProject(filepath.Join(sln.SolutionDir, path))
		if err != nil {
			return sln, err
		}
		sln.ProjectList = append(sln.ProjectList, pro)
	}
	return sln, nil
}

func findAllProject(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return []string{}, err
	}
	re := regexp.MustCompile("[^\"]\"[^\"]+\\.vcxproj\"")
	files := re.FindAllString(string(b), -1)

	var list []string
	for _, v := range files {
		v = strings.Replace(v, "\"", "", -1)
		v = strings.TrimSpace(v)
		list = append(list, v)
	}
	return list, nil
}

func (sln *Sln) CompileCommandsJson(conf string) ([]CompileCommand, error) {
	var cmdList []CompileCommand

	for _, pro := range sln.ProjectList {
		var item CompileCommand

		for _, f := range pro.FindSourceFiles() {
			item.Dir = pro.ProjectDir
			item.File = f

			inc, def, err := pro.FindConfig(conf)
			if err != nil {
				return cmdList, err
			}
			willReplaceEnv := map[string]string{
				"$(SolutionDir)": sln.SolutionDir,
			}
			for k, v := range willReplaceEnv {
				if strings.Contains(inc, k) {
					inc = strings.Replace(inc, k, v, -1)
				}
			}
			def = RemoveBadDefinition(def)
			def = preappend(def, "-D")

			inc = RemoveBadInclude(inc)
			inc = preappend(inc, "-I")

			cmd := "clang-cl.exe " + def + " " + inc + " -c " + f
			item.Cmd = cmd

			cmdList = append(cmdList, item)
		}

	}
	return cmdList, nil
}

func preappend(sepedString string, append string) string {
	defList := strings.Split(sepedString, ";")
	var output string

	for _, v := range defList {
		v = append + v + " "
		output += v
	}
	return output
}
