package types

type State struct {
	TableState        TableState        `json:"table"`
	RelationshipState RelationshipState `json:"relationShip"`
}

type TableState struct {
	Tables  []Table `json:"tables"`
	Indexes []Index `json:"indexes"`
}

type Table struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Comment string   `json:"comment"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Comment  string       `json:"comment"`
	DataType string       `json:"dataType"`
	Default  string       `json:"default"`
	Option   ColumnOption `json:"option"`
	Ui       ColumnUi     `json:"ui"`
}

type ColumnOption struct {
	AutoIncrement bool `json:"autoIncrement"`
	PrimaryKey    bool `json:"primaryKey"`
	Unique        bool `json:"unique"`
	NotNull       bool `json:"notNull"`
}

type ColumnUi struct {
	Pk  bool `json:"pk"`
	Fk  bool `json:"fk"`
	Pfk bool `json:"pfk"`
}

type RelationshipState struct {
	Relationships []Relationship `json:"relationships"`
}

// RelationshipType = ZeroN | OneN | ZeroOne | OneOnly
type Relationship struct {
	Id               string            `json:"id"`
	Identification   bool              `json:"identification"`
	RelationshipType string            `json:"relationshipType"`
	Start            RelationshipPoint `json:"start"`
	End              RelationshipPoint `json:"end"`
	ConstraintName   string            `json:"constraintName?"`
}

type RelationshipPoint struct {
	TableId   string   `json:"tableId"`
	ColumnIds []string `json:"columnIds"`
}

type MemoState struct {
	Memos []Memo `json:"memos"`
}

type Memo struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type Index struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	TableId string        `json:"tableId"`
	Columns []IndexColumn `json:"columns"`
	Unique  bool          `json:"unique"`
}

// OrderType = ASC | DESC
type IndexColumn struct {
	Id        string `json:"id"`
	OrderType string `json:"orderType"`
}

type Node struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Comment string  `json:"comment"`
	Fields  []Field `json:"fields"`
	Edges   []Edge  `json:"edges"`
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
