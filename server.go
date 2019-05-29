package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Server struct {
	Config *Config
}

func newServer(c *Config) *Server {
	return &Server{Config: c}
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
	problem := Problem{Site: q["site"][0], Id: q["id"][0], Config: s.Config}
	b, err := ioutil.ReadFile(problem.submissionPath())
	if err != nil {
		log.Println(err)
	}
	w.Write(b)
}

func (s *Server) launch() {
	http.HandleFunc("/problem", withLog(s.problemHandler))
	http.HandleFunc("/submission", withLog(s.submissionHandler))
	log.Println("Start serving...")
	log.Fatal(http.ListenAndServe(":4567", nil))
}
