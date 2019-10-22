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
		PrimaryKey string
	}

	InsertStmt struct {
		TableName string
		Values []Expr
	}

	SelectStmt struct {
		ColNames []string
		From 	 []string
		Wheres []Expr
	}

	UpdateStmt struct {
 		TableName string
		ColNames []string
 		Set []interface{}
		Where []Expr
	}

	BeginStmt struct {
	}

	CommitStmt struct {
	}

	AbortStmt struct {
	}
)

func (s *CreateTableStmt) stmtNode(){}
func (s *InsertStmt) stmtNode(){}
func (s *SelectStmt) stmtNode(){}
func (s *UpdateStmt) stmtNode(){}
func (s *BeginStmt) stmtNode(){}
func (s *CommitStmt) stmtNode(){}
func (s *AbortStmt) stmtNode(){}

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