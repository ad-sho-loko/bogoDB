package query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCreateTableTokenize(t *testing.T){
	tokenizer := NewTokenizer("create table{int}")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 5)
	assert.Equal(t, tkns[0].kind, CREATE)
	assert.Equal(t, tkns[1].kind, TABLE)
	assert.Equal(t, tkns[2].kind, LBRACE)
	assert.Equal(t, tkns[3].kind, INT)
	assert.Equal(t, tkns[4].kind, RBRACE)
}

func TestLitTokenize(t *testing.T){
	tokenizer := NewTokenizer("a def 1 123")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 4)
	assert.Equal(t, tkns[0].kind, STRING)
	assert.Equal(t, tkns[0].str, "a")
	assert.Equal(t, tkns[1].kind, STRING)
	assert.Equal(t, tkns[1].str, "def")
	assert.Equal(t, tkns[2].kind, NUMBER)
	assert.Equal(t, tkns[2].str, "1")
	assert.Equal(t, tkns[3].kind, NUMBER)
	assert.Equal(t, tkns[3].str, "123")
}

func TestOperatorTokenize(t *testing.T){
	tokenizer := NewTokenizer("{},*=()")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 7)
	assert.Equal(t, tkns[0].kind, LBRACE)
	assert.Equal(t, tkns[1].kind, RBRACE)
	assert.Equal(t, tkns[2].kind, COMMA)
	assert.Equal(t, tkns[3].kind, STAR)
	assert.Equal(t, tkns[4].kind, EQ)
	assert.Equal(t, tkns[5].kind, LPAREN)
	assert.Equal(t, tkns[6].kind, RPAREN)
}

func TestKeywordTokenize(t *testing.T){
	tokenizer := NewTokenizer("create table select where from insert into values update set begin commit rollback")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 13)
	assert.Equal(t, tkns[0].kind, CREATE)
	assert.Equal(t, tkns[1].kind, TABLE)
	assert.Equal(t, tkns[2].kind, SELECT)
	assert.Equal(t, tkns[3].kind, WHERE)
	assert.Equal(t, tkns[4].kind, FROM)
	assert.Equal(t, tkns[5].kind, INSERT)
	assert.Equal(t, tkns[6].kind, INTO)
	assert.Equal(t, tkns[7].kind, VALUES)
	assert.Equal(t, tkns[8].kind, UPDATE)
	assert.Equal(t, tkns[9].kind, SET)
	assert.Equal(t, tkns[10].kind, BEGIN)
	assert.Equal(t, tkns[11].kind, COMMIT)
	assert.Equal(t, tkns[12].kind, ROLLBACK)
}

func TestUpperKeywordTokenize(t *testing.T){
	tokenizer := NewTokenizer("CREATE TABLE")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 2)
	assert.Equal(t, tkns[0].kind, CREATE)
	assert.Equal(t, tkns[1].kind, TABLE)
}

