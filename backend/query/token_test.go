package query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCreateTableTokenize(t *testing.T){
	tokenizer := newTokenizer("create table{int}")
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
	tokenizer := newTokenizer("abc def")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 2)
	assert.Equal(t, tkns[0].kind, STRING)
	assert.Equal(t, tkns[0].str, "abc")
	assert.Equal(t, tkns[1].kind, STRING)
	assert.Equal(t, tkns[1].str, "def")
}

func TestOperatorTokenize(t *testing.T){
	tokenizer := newTokenizer("{},*")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 4)
	assert.Equal(t, tkns[0].kind, LBRACE)
	assert.Equal(t, tkns[1].kind, RBRACE)
	assert.Equal(t, tkns[2].kind, COMMA)
	assert.Equal(t, tkns[3].kind, STAR)
}

func TestKeywordTokenize(t *testing.T){
	tokenizer := newTokenizer("create table select where from")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 5)
	assert.Equal(t, tkns[0].kind, CREATE)
	assert.Equal(t, tkns[1].kind, TABLE)
	assert.Equal(t, tkns[2].kind, SELECT)
	assert.Equal(t, tkns[3].kind, WHERE)
	assert.Equal(t, tkns[4].kind, FROM)
}

func TestUpperKeywordTokenize(t *testing.T){
	tokenizer := newTokenizer("CREATE TABLE")
	tkns, err := tokenizer.Tokenize()

	if err != nil{
		log.Fatal(err)
	}

	assert.Equal(t, len(tkns), 2)
	assert.Equal(t, tkns[0].kind, CREATE)
	assert.Equal(t, tkns[1].kind, TABLE)
}

