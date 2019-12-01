package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Submission struct {
	Problem Problem
	Path    string
	Ext     string
}

func parsePath(path string, config *Config) (string, string, string, error) {
	e := filepath.Join(
		"("+config.SolutionDir+")",
		"([^"+string(filepath.Separator)+"]+)",
		"([^."+string(filepath.Separator)+"]+)",
	)
	re := regexp.MustCompile(e)
	ms := re.FindSubmatch([]byte(path))
	if len(ms) != 4 {
		return "", "", "", errors.New("Could not handle path \"" + path + "\"")
	}
	site := string(ms[2])
	id := string(ms[3])
	var ext string
	if string(ms[1]) == config.SolutionDir {
		ext = filepath.Ext(path)[1:]
	} else {
		ext = config.DefaultLanguage
	}
	return site, id, ext, nil
}

func newSubmission(name string, config *Config) (*Submission, error) {
	site, id, ext, err := parsePath(name, config)
	if err != nil {
		return nil, err
	}
	p, err := loadProblem(site, id, config)
	if err != nil {
		return nil, err
	}
	return &Submission{
		Path:    filepath.Join(config.SolutionDir, site, id+"."+ext),
		Ext:     ext,
		Problem: *p}, nil
}

func (s *Submission) writePreprocessedSource(config *Config, source []byte) (string, error) {
	rel, err := filepath.Rel(config.SolutionDir, s.Path)
	if err != nil {
		return "", err
	}
	sourcePath := filepath.Join(config.SolutionDir, "tmp", rel)
	err = os.MkdirAll(filepath.Dir(sourcePath), 0755)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(sourcePath, source, 0644)
	if err != nil {
		return "", err
	}
	return sourcePath, nil
}

func (s *Submission) test(config *Config) *TestResult {
	var ex string
	var args []string
	var err error
	source, err := s.preprocess(config)
	if err != nil {
		return newFailedTestResult(CE, err)
	}

	sourcePath, err := s.writePreprocessedSource(config, source)
	if err != nil {
		return newFailedTestResult(IE, err)
	}

	if config.Languages[s.Ext].Compile != "" {
		ex, err = s.compile(config, sourcePath)
		if ex != "" {
			defer os.Remove(ex)
		}
		if err != nil {
			return newFailedTestResult(CE, err)
		}
	} else if config.Languages[s.Ext].Interpret != "" {
		ex = config.Languages[s.Ext].Interpret
		args = append(args, sourcePath)
	} else {
		return newFailedTestResult(CE, errors.New("language config of "+s.Ext+" must have `interpret` or `compile`"))
	}
	return newTestResult(sourcePath, s, s.Problem.test(ex, args))
}

func (s *Submission) shouldPreprocess(config *Config) bool {
	return config.Languages[s.Ext].Preprocess != ""
}

func (s *Submission) preprocess(config *Config) ([]byte, error) {
	var res []byte
	if !s.shouldPreprocess(config) {
		f, err := os.Open(s.Path)
		defer f.Close()
		_, err = f.Read(res)
		if err != nil {
			return res, err
		}
		return res, nil
	}

	out, err := exec.Command(config.Languages[s.Ext].Preprocess, s.Path).CombinedOutput()
	if err != nil {
		return res, errors.New(err.Error() + "\n" + string(out))
	}
	return out, nil
}

func (s *Submission) compile(config *Config, path string) (string, error) {
	log.Println("compilation start")
	tmpfile, err := ioutil.TempFile(".", "")
	if err != nil {
		return "", err
	}
	tmpName := "." + string(filepath.Separator) + tmpfile.Name()
	args := strings.Fields(config.Languages[s.Ext].Compile)
	log.Println("\x1b[1m" + strings.Join(args, " ") + " TMP_FILE " + s.Path + "\x1b[0m")
	args = append(args, tmpfile.Name(), path)
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return tmpName, errors.New(err.Error() + "\n" + string(out))
	}
	log.Println("compilation done")
	return tmpName, nil
}
