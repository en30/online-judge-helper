package main

import (
	"errors"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ProblemKey struct {
	Site string
	Id   string
}

type Server struct {
	Config      *Config
	Connections map[ProblemKey]*websocket.Conn
}

func newServer(c *Config) *Server {
	return &Server{Config: c, Connections: make(map[ProblemKey]*websocket.Conn)}
}

func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
		log.Printf("%v %v %v", r.Method, r.URL.Path, r.Proto)
	}
}

func (s *Server) openEditor(problem *Problem) error {
	editor := s.Config.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	args := strings.Split(editor, " ")
	args = append(args, problem.submissionPath())
	return exec.Command(args[0], args[1:]...).Start()
}

func (s *Server) problemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var problem Problem
	err := problem.parse(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	problem.Config = s.Config

	err = problem.save()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	os.MkdirAll(problem.submissionDir(), 0755)
	s.openEditor(&problem)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) submissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	q := r.URL.Query()
	key := ProblemKey{Site: q["site"][0], Id: q["id"][0]}
	websocket.Handler(s.watchSubmissionHandler(key)).ServeHTTP(w, r)
}

func (s *Server) watchSubmissionHandler(key ProblemKey) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		s.Connections[key] = ws
		problem, err := loadProblem(key.Site, key.Id, s.Config)
		if err != nil {
			log.Fatal(err)
		}
		sub, err := newSubmission(problem.submissionPath(), s.Config)
		if err != nil {
			log.Println(err)
			return
		}
		source, err := sub.preprocess(s.Config)
		if err != nil {
			log.Println(err)
			return
		}
		ws.Write(source)
		io.Copy(ioutil.Discard, ws)
	}
}

func (s *Server) sendSubmission(sub *Submission, config *Config) error {
	key := ProblemKey{Site: sub.Problem.Site, Id: sub.Problem.Id}
	if s.Connections[key] == nil {
		return errors.New("No connection for " + sub.Problem.Site + "/" + sub.Problem.Id)
	}
	source, err := sub.preprocess(config)
	if err != nil {
		return err
	}
	s.Connections[key].Write(source)
	return nil
}

func (s *Server) launch() {
	http.HandleFunc("/problem", withLog(s.problemHandler))
	http.HandleFunc("/submission", withLog(s.submissionHandler))
	log.Println("Start serving...")
	log.Fatal(http.ListenAndServe(":4567", nil))
}
