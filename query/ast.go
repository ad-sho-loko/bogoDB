package query

type Stmt interface {
	stmtNode()
}

// statements
type (
	CreateTableStmt struct {
		TableName string
		ColNames []string
		ColTypes []string
	}

	InsertStmt struct {
		TableName string
		Values []Expr
	}

	SelectStmt struct {
		ColNames []string
		From *From
		Wheres []Expr
	}

	UpdateStmt struct {
 		TableName string
		ColNames []string
		Set []Expr
	}
)

func (s *CreateTableStmt) stmtNode(){}
func (s *InsertStmt) stmtNode(){}
func (s *SelectStmt) stmtNode(){}
func (s *UpdateStmt) stmtNode(){}

// clauses
type From struct {
	TableNames []string
}

type Where struct {
	Cond []Expr
}

// expressions
type Expr interface {
	exprNode()
}

type Eq struct {
	left Expr
	right Expr
}

type Lit struct {
	v string
}

func (l *Eq) exprNode(){}
func (l *Lit) exprNode(){}