package query

import (
	"fmt"
	"github.com/pkg/errors"
)

// Parser parses input tokens.
type Parser struct {
	tokens []*Token
	pos int
	errors []error
}

func NewParser(tokens []*Token) *Parser{
	return &Parser{
		tokens:tokens,
	}
}

func (p *Parser) consume(kind TokenKind) bool{
	if p.pos < len(p.tokens) && p.tokens[p.pos].kind == kind{
		p.pos++
		return true
	}
	return false
}

func (p *Parser) errorExpected(e error){
	p.errors = append(p.errors, e)
}

func (p *Parser) expectOr(kinds ...TokenKind) *Token{
	for _, k := range kinds{
		if p.tokens[p.pos].kind == k{
			tkn := p.tokens[p.pos]
			p.pos++
			return tkn
		}
	}
	p.errorExpected(fmt.Errorf("expected %s or %s, but %s", kinds[0], kinds[1], p.tokens[p.pos].kind))
	return nil
}

func (p *Parser) expect(kind TokenKind) *Token{
	if p.tokens[p.pos].kind == kind{
		tkn := p.tokens[p.pos]
		p.pos++
		return tkn
	}

	p.errorExpected(fmt.Errorf("expected %s, but %s", kind, p.tokens[p.pos].kind))
	return nil
}

func (p *Parser) expr() Expr {
	token := p.tokens[p.pos]

	if p.consume(NUMBER) || p.consume(STRING){
		return &Lit{v:token.str}
	}

	p.errorExpected(fmt.Errorf("expr failed"))
	return nil
}

func (p *Parser) eq() Expr{
	left := p.expr()

	if p.consume(EQ){
		right := p.expr()
		return &Eq{
			left:left,
			right:right,
		}
	}
	return left
}

func (p *Parser) fromClause() []string{
	s := p.expect(STRING)
	return []string{s.str}
}

func (p *Parser) whereClause() []Expr{
	var exprs []Expr
	exprs = append(exprs, p.eq())
	return exprs
}

func (p *Parser) selectStmt() Stmt{
	// select
	tkn := p.expectOr(STAR, STRING)

	selectNode := &SelectStmt{
		ColNames:[]string{tkn.str},
	}

	// from
	if p.consume(FROM){
		from := p.fromClause()
		selectNode.From = from
	}

	// where
	if p.consume(WHERE){
		selectNode.Wheres = p.whereClause()
	}

	return selectNode
}

func (p *Parser) updateTableStmt() Stmt{
	tblName := p.expect(STRING)
	p.expect(SET)

	var cols []string
	var sets []interface{}

	for{
		col := p.expect(STRING)
		cols = append(cols, col.str)

		p.expect(EQ)

		set := p.expectOr(STRING ,NUMBER)
		sets = append(sets, set.str)

		if !p.consume(COMMA){
			break
		}
	}

	// where
	var exprs []Expr
	if p.consume(WHERE){
		exprs = p.whereClause()
	}

	return &UpdateStmt{
		TableName:tblName.str,
		ColNames:cols,
		Set:sets,
		Where:exprs,
	}
}

func (p *Parser) insertTableStmt() Stmt{
	p.expect(INTO)
	tblName := p.expect(STRING)
	p.expect(VALUES)
	p.expect(LPAREN)

	var exprs []Expr
	for {
		exprs = append(exprs, p.eq())
		if p.consume(RPAREN){
			break
		}
		p.expect(COMMA)
	}

	return &InsertStmt{
		TableName:tblName.str,
		Values:exprs,
	}
}

func (p *Parser) createTableStmt() Stmt{
	p.expect(TABLE)
	tblName := p.expect(STRING)
	p.expect(LBRACE)

	var colNames []string
	var colTypes []string
	var pk string

	for{
		colName := p.expect(STRING)
		p.expect(INT)
		colNames = append(colNames, colName.str)
		colTypes = append(colTypes, "int")
		if p.consume(PRIMARY){
			p.expect(KEY)
			pk = colName.str
		}

		if !p.consume(COMMA){
			break
		}
	}

	p.expect(RBRACE)

	return &CreateTableStmt{
		TableName:tblName.str,
		ColNames:colNames,
		ColTypes:colTypes,
		PrimaryKey:pk,
	}
}

func (p *Parser) Parse() (Stmt, []error){
	// create table
	if p.consume(CREATE){
		return p.createTableStmt(), p.errors
	}

	// select
	if p.consume(SELECT){
		return p.selectStmt(), p.errors
	}

	// insert
	if p.consume(INSERT){
		return p.insertTableStmt(), p.errors
	}

	// update
	if p.consume(UPDATE){
		return p.updateTableStmt(), p.errors
	}

	if p.consume(BEGIN){
		return &BeginStmt{}, nil
	}

	if p.consume(COMMIT){
		return &CommitStmt{}, nil
	}

	if p.consume(ROLLBACK){
		return &AbortStmt{}, nil
	}

	// unexpected query comes
	p.errors = append(p.errors, errors.New("unexpected query"))
	return nil, p.errors
}