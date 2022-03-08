package types

type Node struct {
	Name   string
	Fields []Field
	Edges  []Edge
}

type Edge struct {
	Name      string
	Type      string // 1..N | 0..1 | 0..N | 1..1
	Direction string // In | Out
}

type Field struct {
	Name          string
	Type          string
	Comment       string
	Default       string
	Pk            bool
	Fk            bool
	Pfk           bool
	Unique        bool
	AutoIncrement bool
	Sensitive     bool
}

type File struct {
	Path   string
	Buffer string
}
