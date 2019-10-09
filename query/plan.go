package query

import (
	"github.com/ad-sho-loko/bogodb/storage"
	"github.com/pkg/errors"
)

// Planner is to plan
type Planner struct {
	q Query
}

type Plan struct {
	scanners Scanner
}

type Scanner interface {
	Scan(storage *storage.Storage) []*storage.Tuple
}

// scanner
type(
	SeqScan struct{
		tblName string
	}

	IndexScan struct {
		indexName string
	}
)

func NewPlanner(q Query) *Planner{
	return &Planner{
		q:q,
	}
}

func (p *Planner) planSelect(q *SelectQuery) (*Plan, error){
	return &Plan{
		scanners:&SeqScan{
			tblName:q.From[0].Name,
		},
	}, nil
}

func (p *Planner) PlanMain() (*Plan, error){
	switch concrete := p.q.(type) {
	case *SelectQuery:
		return p.planSelect(concrete)
	case *CreateTableQuery:
		// do nothing
		return nil, nil
	}
	return nil, errors.New("unexpected query when planning")
}