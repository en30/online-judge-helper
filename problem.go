package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Restriction struct {
	TimeLimit time.Duration
}

func (r *Restriction) save(path string) error {
	b, err := json.Marshal(*r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

type Problem struct {
	Site        string      `json:",omitempty"`
	Id          string      `json:",omitempty"`
	Restriction Restriction `json:",omitempty"`
	TestCases   []TestCase  `json:",omitempty"`
	Config      *Config     `json:"-"`
}

func loadProblem(site string, id string, config *Config) (*Problem, error) {
	p := Problem{Site: site, Id: id, Config: config}
	p.loadRestriction()
	p.loadTestCases()
	return &p, nil
}

func (p *Problem) loadRestriction() error {
	b, err := ioutil.ReadFile(p.restrictionPath())
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &p.Restriction)
}

func (p *Problem) loadTestCases() error {
	matches, err := filepath.Glob(filepath.Join(p.testDir(), "*.txt"))
	if err != nil {
		return err
	}
	idx := map[string]int{}
	for _, m := range matches {
		title := filepath.Base(m)
		title = strings.TrimSuffix(title, ".in.txt")
		title = strings.TrimSuffix(title, ".out.txt")
		if i, ok := idx[title]; !ok {
			tc := TestCase{Title: title}
			i = len(p.TestCases)
			p.TestCases = append(p.TestCases, tc)
			idx[title] = i
		}
		tc := &p.TestCases[idx[title]]
		b, err := ioutil.ReadFile(m)
		if err != nil {
			return err
		}
		if strings.HasSuffix(m, ".in.txt") {
			tc.Input = string(b)
		} else if strings.HasSuffix(m, "out.txt") {
			tc.Output = string(b)
		}
	}
	return nil
}

func (p *Problem) parse(r io.ReadCloser) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}
	p.Site = filepath.Base(p.Site)
	if p.Site == "" || p.Site == ".." || p.Site == "." {
		p.Site = "unknown"
	}
	return nil
}

func (p *Problem) save() error {
	os.MkdirAll(p.testDir(), 0755)
	p.Restriction.save(p.restrictionPath())
	for _, tc := range p.TestCases {
		err := tc.save(p.testDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Problem) restrictionPath() string {
	return filepath.Join(p.testDir(), "restriction.json")
}

func (p *Problem) testDir() string {
	return filepath.Join(p.Config.TestDir, p.Site, p.Id)
}

func (p *Problem) submissionDir() string {
	return filepath.Join(p.Config.SolutionDir, p.Site)
}

func (p *Problem) submissionPath() string {
	return filepath.Join(p.submissionDir(), p.Id+"."+p.Config.DefaultLanguage)
}

func (p *Problem) test(command string, args []string) []TestCaseResult {
	var wg sync.WaitGroup
	results := make([]TestCaseResult, len(p.TestCases))
	for i, _ := range p.TestCases {
		wg.Add(1)
		go func(i int) {
			r := p.TestCases[i].test(command, args, p.Restriction.TimeLimit)
			results[i] = *r
			wg.Done()
		}(i)
	}
	wg.Wait()
	return results
}

type TestResult struct {
	Submission Submission
	Status     Status
	Results    []TestCaseResult
	Error      error
	SourcePath string
}

func newFailedTestResult(s Status, err error) *TestResult {
	return &TestResult{Status: s, Error: err}
}

func newTestResult(sourcePath string, sub *Submission, results []TestCaseResult) *TestResult {
	s := AC
	for _, r := range results {
		if s == AC && r.Status != AC {
			s = r.Status
		}
	}
	return &TestResult{Submission: *sub, Status: s, Results: results, SourcePath: sourcePath}
}
