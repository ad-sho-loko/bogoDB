package query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

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

	p := NewParser()
	node, err := p.parse(tokens)
	if err != nil{
		log.Fatal(err)
	}

	n := node.(*CreateTableNode)
	assert.Equal(t, "users", n.TableName)
	assert.Equal(t, "id", n.ColNames[0])
	assert.Equal(t, "int", n.ColTypes[0])
}

func TestParseSelect(t *testing.T){
	tokens := []*Token{
		{kind:SELECT},
		// {kind:STAR},
		{kind:STRING, str:"id"},
		{kind:FROM},
		{kind:STRING, str:"users"},
	}

	p := NewParser()
	node, err := p.parse(tokens)
	if err != nil{
		log.Fatal(err)
	}

	n := node.(*SelectNode)
	assert.Equal(t, "users", n.From.TableNames[0])
	assert.Equal(t, "id", n.ColNames[0])
}
