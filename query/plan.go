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
		tblName string
		indexName string
	}
)

func NewPlanner(q Query) *Planner{
	return &Planner{
		q:q,
	}
}

func (p *Planner) planSelect(q *SelectQuery) (*Plan, error){
	// if where contains a primary key, use index scan.
	for _, w := range q.Where{
		eq, ok := w.(*Eq)
		if !ok{
			continue
		}
		col, ok := eq.left.(*Lit)
		if !ok{
			continue
		}
		for _, c := range q.Cols{
			if col.v == c.Name && c.Primary{
				return &Plan{
					scanners:&IndexScan{
						tblName:q.From[0].Name,
						indexName:col.v,
					},
				}, nil
			}
		}
	}

	// use seqscan
	return &Plan{
		scanners:&SeqScan{
			// FIXME
			tblName:q.From[0].Name,
		},
	}, nil
}

func (p *Planner) PlanMain() (*Plan, error){
	switch concrete := p.q.(type) {
	case *SelectQuery:
		return p.planSelect(concrete)
	case *CreateTableQuery, *InsertQuery, *UpdateQuery:
		return nil, nil
	case *BeginQuery, *CommitQuery, *AbortQuery:
		return nil, nil
	}
	return nil, errors.New("unexpected query when planning")
}