package query

import (
	"fmt"
	"github.com/ad-sho-loko/bogodb/meta"
	"github.com/ad-sho-loko/bogodb/storage"
	"strconv"
)

// Analyzer analyze the parsed sql.
// Roles:
//   - Fetch an actual scheme from the table name.
//   - Validate the rules of sql.
type Analyzer struct {
	catalog *storage.Catalog
}

type Query interface {
	evalQuery()
}

type SelectQuery struct {
	Cols  []*meta.Column
	From  []*meta.Table
	Where []Expr
}

type CreateTableQuery struct {
	Scheme *meta.Scheme
}

type UpdateQuery struct {
	Table  *meta.Table
	Cols []*meta.Column
	Set []interface{}
}

type InsertQuery struct {
	Table  *meta.Table
	Values []interface{}
}

func (q *SelectQuery) evalQuery(){}
func (q *InsertQuery) evalQuery(){}
func (q *CreateTableQuery) evalQuery(){}
func (q *UpdateQuery) evalQuery(){}

func NewAnalyzer(catalog *storage.Catalog) *Analyzer{
	return &Analyzer{
		catalog:catalog,
	}
}

func (a *Analyzer) analyzeInsert(n *InsertStmt) (*InsertQuery, error){
	var q InsertQuery

	// analyze `into`
	if !a.catalog.HasScheme(n.TableName){
		return nil, fmt.Errorf("insert failed : `%s` doesn't exists", n.TableName)
	}
	scheme := a.catalog.FetchScheme(n.TableName)

	t := &meta.Table{
		Name:n.TableName,
	}

	// analyze `values`
	if len(n.Values) != len(scheme.ColNames){
		return nil, fmt.Errorf("insert failed : `values` should be same length")
	}

	var lits []string
	for _, l := range n.Values{
		num := l.(*Lit)
		lits = append(lits, num.v)
	}

	// FIXME
 	var values []interface{}
	for i, v := range lits{
		if scheme.ColTypes[i] == meta.Int{
			n, _ := strconv.Atoi(v)
			values = append(values, n)
		}else if scheme.ColTypes[i] == meta.Varchar{
			values = append(values, v)
		}else{
			return nil, fmt.Errorf("insert failed : unexpected types parsed")
		}
	}

	q.Table = t
	q.Values = values
	return &q, nil
}

func (a *Analyzer) analyzeSelect(n *SelectStmt) (*SelectQuery, error){
	var q SelectQuery

	// analyze `from`
	var schemes []*meta.Scheme
	for _, name := range n.From.TableNames{
		scheme := a.catalog.FetchScheme(name)
		if scheme == nil{
			return nil, fmt.Errorf("select failed :table `%s` doesn't exist", name)
		}
		schemes = append(schemes, scheme)
	}

	// analyze `select`
	var cols []*meta.Column
	for _, colName := range n.ColNames{
		found := false
		for _, scheme := range schemes{
			for _, col := range scheme.ColNames{
				if col == colName{
					found = true
					col := &meta.Column{
						Name:colName,
					}
					cols = append(cols, col)
				}
			}
		}

		if !found{
			return nil, fmt.Errorf("select failed : column `%s` doesn't exist", colName)
		}
	}

	// analyze `where`

	/*
	for _, e := range n.Wheres{
		eq, ok := e.(*Eq)
		if !ok{
		}

		// left must be col-name.
		colName := eq.left.(*Lit).v
		for _, col := range cols{
			if col.Name == colName{

			}
		}
	}
	*/

	var tables []*meta.Table
	for _, s := range schemes{
		table := s.ConvertTable()
		tables = append(tables, table)
	}

	q.From = tables
	q.Cols = cols
	q.Where = n.Wheres
	return &q, nil
}

func (a *Analyzer) analyzeUpdate(n *UpdateStmt) (*UpdateQuery, error) {
	var q UpdateQuery

	// analyze `update`
	if !a.catalog.HasScheme(n.TableName){
		return nil, fmt.Errorf("insert failed : `%s` doesn't exists", n.TableName)
	}
	t := &meta.Table{
		Name:n.TableName,
	}

	// analyze `set`
	// scheme := a.catalog.FetchScheme(n.TableName)
	/*
	var values []interface{}
	for _, v := range n.Set{
		return nil, fmt.Errorf("insert failed : unexpected types parsed")
	}
	*/

	q.Table = t
	return &q, nil
}

func (a *Analyzer) analyzeCreateTable(n *CreateTableStmt) (*CreateTableQuery, error){
	var q CreateTableQuery

	if a.catalog.HasScheme(n.TableName){
		return nil, fmt.Errorf("create table failed : table name `%s` already exists", n.TableName)
	}

	var types []meta.ColType
	for _, typ := range n.ColTypes{
		if typ == "int"{
			types = append(types, meta.Int)
		} else if typ == "varchar"{
			types = append(types, meta.Varchar)
		}
	}

	q.Scheme = meta.NewScheme(n.TableName, n.ColNames, types)
	return &q, nil
}

func (a *Analyzer) AnalyzeMain(stmt Stmt) (Query, error){
	switch concrete := stmt.(type) {
	case *SelectStmt:
		return a.analyzeSelect(concrete)
	case *CreateTableStmt:
		return a.analyzeCreateTable(concrete)
	case *InsertStmt:
		return a.analyzeInsert(concrete)
	case *UpdateStmt:
		return a.analyzeUpdate(concrete)
	}

	return nil, fmt.Errorf("failed to analyze query")
}