package query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestEq(t *testing.T){
	tokens := []*Token{
		{kind:STRING, str:"1"},
		{kind:EQ},
		{kind:STRING, str:"2"},
	}

	p := NewParser(tokens)
	node := p.eq().(*Eq)
	left := node.left.(*Lit)
	right := node.right.(*Lit)

	assert.Equal(t, "1", left.v)
	assert.Equal(t, "2", right.v)
}


func TestParseCreateTable(t *testing.T){
	tokens := []*Token{
		{kind:CREATE},
		{kind:TABLE},
		{kind:STRING, str:"users"},
		{kind:LBRACE},
		{kind:STRING, str:"id"},
		{kind:INT, str:"int"},
		{kind:RBRACE},
	}

	p := NewParser(tokens)
	node, err := p.Parse()
	if err != nil{
		log.Fatal(err)
	}

	n := node.(*CreateTableStmt)
	assert.Equal(t, "users", n.TableName)
	assert.Equal(t, "id", n.ColNames[0])
	assert.Equal(t, "int", n.ColTypes[0])
}

func TestParseSelect(t *testing.T){
	tokens := []*Token{
		{kind:SELECT},
		{kind:STRING, str:"id"},
		{kind:FROM},
		{kind:STRING, str:"users"},
		{kind:WHERE, str:"where"},
		{kind:STRING, str:"1"},
		{kind:EQ},
		{kind:STRING, str:"2"},
	}

	p := NewParser(tokens)
	node, err := p.Parse()
	if err != nil{
		log.Fatal(err)
	}

	n := node.(*SelectStmt)
	assert.Equal(t, "users", n.From.TableNames[0])
	assert.Equal(t, "id", n.ColNames[0])
}


func TestParseInsert(t *testing.T){
	tokens := []*Token{
		{kind:INSERT},
		{kind:INTO},
		{kind:STRING, str:"users"},
		{kind:VALUES},
		{kind:LPAREN},
		{kind:NUMBER, str:"1"},
		{kind:RPAREN},
	}

	p := NewParser(tokens)
	node, err := p.Parse()
	if err != nil{
		log.Fatal(err)
	}

	n := node.(*InsertStmt)
	assert.Equal(t, "users", n.TableName)
	assert.Equal(t, "1", n.Values[0].(*Lit).v)
}