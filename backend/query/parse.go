package query

import (
	"fmt"
	"github.com/go-ffmt/ffmt"
	"github.com/pkg/errors"
)

type Parser struct {
	tokens []*Token
	pos int
}

func NewParser() *Parser{
	return &Parser{
	}
}

type Node interface {
	ErrInfo() string
}

type CreateTableNode struct {
	TableName string
	ColNames []string
	ColTypes []string
}

func (c *CreateTableNode) ErrInfo() string{
	return ffmt.Sputs(c)
}

type SelectNode struct {
	ColNames []string
	From *FromNode
	Wheres Node
}

func (c *SelectNode) ErrInfo() string{
	return ffmt.Sputs(c)
}

type FromNode struct {
	TableNames []string
}

type WhereNode struct {
}

func (p *Parser) consume(kind TokenKind) bool{
	if p.tokens[p.pos].kind == kind{
		p.pos++
		return true
	}
	return false
}

func (p *Parser) expectOr(kinds ...TokenKind) (*Token, error){
	for _, k := range kinds{
		if p.tokens[p.pos].kind == k{
			tkn := p.tokens[p.pos]
			p.pos++
			return tkn, nil
		}
	}
	return nil, fmt.Errorf("expected %s or %s, but %s", kinds[0], kinds[1], p.tokens[p.pos].kind)
}

func (p *Parser) expect(kind TokenKind) (*Token, error){
	if p.tokens[p.pos].kind == kind{
		tkn := p.tokens[p.pos]
		p.pos++
		return tkn, nil
	}
	return nil, fmt.Errorf("expected %s, but %s", kind, p.tokens[p.pos].kind)
}

func (p *Parser) fromPhase() (*FromNode, error){
	s, err := p.expect(STRING)
	if err != nil{
		return nil, err
	}

	return &FromNode{
		TableNames:[]string{s.str},
	}, nil
}

func (p *Parser) selectStmt() (Node, error){
	// select
	tkn, err := p.expectOr(STAR, STRING)
	if err != nil{
		return nil, err
	}
	selectNode := &SelectNode{
		ColNames:[]string{tkn.str},
	}

	// from
	if p.consume(FROM){
		from, err := p.fromPhase()
		if err != nil{
			return nil, err
		}
		selectNode.From = from
	}

	return selectNode, nil
}

func (p *Parser) createTableStmt() (Node, error){
	if _, err := p.expect(TABLE); err != nil { return nil, err }
	tblName, err := p.expect(STRING); if err != nil { return nil, err }
	if _, err := p.expect(LBRACE);  err != nil { return nil, err }

	var colNames []string
	var colTypes []string

	for{
		colName, err := p.expect(STRING); if err != nil { return nil, err }
		colType, err := p.expect(INT); if err != nil { return nil, err }
		colNames = append(colNames, colName.str)
		colTypes = append(colTypes, colType.str)
		if !p.consume(COMMA){
			break
		}
	}

	if _, err := p.expect(RBRACE);  err != nil { return nil, err }

	return &CreateTableNode{
		TableName:tblName.str,
		ColNames:colNames,
		ColTypes:colTypes,
	}, nil
}

func (p *Parser) top() (Node, error){
	// create table
	if p.consume(CREATE){
		return p.createTableStmt()
	}

	if p.consume(SELECT){
		return p.selectStmt()
	}

	return nil, errors.New("failed to parse query")
}

func (p *Parser) parse(tokens []*Token) (Node, error){
	p.tokens = tokens
	return p.top()
}

func (p *Parser) ParseMain(query string) (Node, error){
	tokenizer := newTokenizer(query)
	tokens, err := tokenizer.Tokenize()

	if err != nil{
		return nil, err
	}

	return p.parse(tokens)
}
