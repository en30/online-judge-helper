package main

import (
	"container/ring"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func watchRecursive(basePath string, watcher *fsnotify.Watcher) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		if filepath.HasPrefix(rel, ".git") {
			return nil
		}
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		return nil
	})
}

type EventLog struct {
	Name    string
	Time    time.Time
	ModTime time.Time
}

func newEventLog(event fsnotify.Event) (*EventLog, error) {
	s, err := os.Stat(event.Name)
	if err != nil {
		return nil, err
	}

	return &EventLog{Name: event.Name, Time: time.Now(), ModTime: s.ModTime()}, nil
}

type Throttler struct {
	events   *ring.Ring
	duration time.Duration
}

func newThrottler(n int, duration time.Duration) *Throttler {
	return &Throttler{ring.New(n), duration}
}

func (t *Throttler) tooFrequent(el *EventLog) bool {
	n := t.events.Len()
	r := t.events
	now := time.Now()
	for i := 0; i < n; i++ {
		if r.Value == nil {
			return false
		}
		e := r.Value.(*EventLog)
		if e.Name == el.Name && (e.ModTime == el.ModTime || now.Sub(e.Time) < t.duration) {
			return true
		}
		r = r.Prev()
	}
	return false
}

func (t *Throttler) add(e *EventLog) {
	t.events = t.events.Next()
	t.events.Value = e
}

func watch(sc chan *Submission, config *Config) {
	tmpl := template.New("test_result.tmpl")
	tmpl, err := tmpl.Funcs(template.FuncMap{
		"cage": func(title string, v string) string {
			return "\u250f \x1b[1m" + title + "\x1b[0m\n\u2503 " + strings.Replace(v, "\n", "\n\u2503 ", -1) + "\n\u2517"
		},
		"leftpad": func(pad string, v string) string {
			return pad + strings.Replace(v, "\n", "\n"+pad, -1)
		},
		"bold": func(s string) string {
			return "\x1b[1m" + s + "\x1b[0m"
		},
		"red": func(s string) string {
			return "\x1b[31m" + s + "\x1b[0m"
		},
		"green": func(s string) string {
			return "\x1b[32m" + s + "\x1b[0m"
		},
		"yellow": func(s string) string {
			return "\x1b[33m" + s + "\x1b[0m"
		},
		"blue": func(s string) string {
			return "\x1b[34m" + s + "\x1b[0m"
		},
	}).ParseFiles("test_result.tmpl", "status.tmpl", "test_case_result.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Println("Start watching...")
	err = watchRecursive(config.SolutionDir, watcher)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("watching", filepath.Join(config.SolutionDir, "**", "*"))

	t := newThrottler(8, 50*time.Millisecond)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if strings.Index(path.Base(event.Name), "#") != -1 {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				el, err := newEventLog(event)
				if err != nil {
					log.Fatal(err)
				}

				if t.tooFrequent(el) {
					continue
				}
				t.add(el)

				log.Println("test triggered by", event.Name, time.Now())
				var res *TestResult
				s, err := newSubmission(event.Name, config)
				if err != nil {
					res = newFailedTestResult(IE, err)
				} else {
					sc <- s
					res = s.test(config)
				}
				err = tmpl.Execute(os.Stdout, res)
				if err != nil {
					log.Fatal(err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
