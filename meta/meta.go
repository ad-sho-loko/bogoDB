package meta

// Table is the definition of tables and Columns
type Scheme struct {
	TblName    string    `json:"tblName"`
	ColNames   []string  `json:"colNames"`
	ColTypes   []ColType `json:"colTypes"`
	PrimaryKey string    `json:"pk"`
}

func NewScheme(tblName string, colNames []string, colTypes []ColType, pk string) *Scheme {
	return &Scheme{
		TblName:    tblName,
		ColNames:   colNames,
		ColTypes:   colTypes,
		PrimaryKey: pk,
	}
}

func (s *Scheme) ConvertTable() *Table {
	var t Table
	t.Name = s.TblName

	var columns []Column
	for i := range s.ColNames {
		var col Column
		col.Name = s.ColNames[i]
		col.Type = string(s.ColTypes[i])
		columns = append(columns, col)
	}

	t.Columns = columns
	return &t
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name    string
	Type    string
	Primary bool
}

type ColType uint8

const (
	Int ColType = iota
	Varchar
)

func (c ColType) String() string {
	if c == Int {
		return "int"
	}

	if c == Varchar {
		return "varchar"
	}

	return "undefined"
}

type ResultSet struct {
	Message string
	ColNames []string
	Values []string
}

func NewWithMessage(message string) *ResultSet{
	return &ResultSet{
		Message:message,
	}
}