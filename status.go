package main

type Status int

const (
	AC Status = iota
	CE
	MLE
	TLE
	RE
	OLE
	IE
	WA
)

func (s Status) Success() bool {
	return s == AC
}
