package types

type Node struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
	Edges  []Edge  `json:"edges"`
}

type Edge struct {
	ID        string `json:"id"`
	Field     Field  `json:"field"`
	Reference Field  `json:"reference"`
	Name      string `json:"name"`
	Type      string `json:"type"`      // 1..N | 0..1 | 0..N | 1..1
	Direction string `json:"direction"` // In | Out
}

type Field struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Comment       string `json:"cemment"`
	Default       string `json:"default"`
	Pk            bool   `json:"pk"`
	Fk            bool   `json:"fk"`
	Pfk           bool   `json:"pfk"`
	Unique        bool   `json:"unique"`
	AutoIncrement bool   `json:"autoIncrement"`
	Sensitive     bool   `json:"sensitive"`
	Nullable      bool   `json:"nullable"`
}

type File struct {
	Path   string `json:"path"`
	Buffer string `json:"buffer"`
}
