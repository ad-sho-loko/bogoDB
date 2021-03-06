package db

import (
	"github.com/ad-sho-loko/bogodb/meta"
	"github.com/ad-sho-loko/bogodb/query"
	"github.com/ad-sho-loko/bogodb/storage"
	"log"
	"os"
	"os/signal"
	"strings"
)

type BogoDb struct {
	exit        chan int
	contexts    map[string]*dbSession
	storage     *storage.Storage
	catalog     *storage.Catalog
	tranManager *storage.TransactionManager
	home        string
}

type dbSession struct {
	tran *storage.Transaction
}

func NewBogoDb() (*BogoDb, error) {
	// load the catalog if exists
	home, ok := os.LookupEnv("BOGO_HOME")
	if !ok {
		// default
		home = ".bogo/"
		if _, err := os.Stat(home); os.IsNotExist(err) {
			err := os.Mkdir(home, 0777)
			if err != nil {
				panic(err)
			}
		}
	}

	catalog, err := storage.LoadCatalog(home)
	if err != nil {
		return nil, err
	}

	return &BogoDb{
		catalog:     catalog,
		storage:     storage.NewStorage(home),
		tranManager: storage.NewTransactionManager(),
		contexts:    make(map[string]*dbSession),
		exit:        make(chan int, 1),
		home:        home,
	}, nil
}

func (db *BogoDb) Init() {
	// set up the signal handler
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		for {
			r := <-sig
			if r == os.Interrupt {
				log.Println("ctrl+c detected, shutdown soon....")
				db.Terminate()
			}
		}
	}()
}

func (db *BogoDb) Execute(q string, userAgent string) (string, error) {
	var trn *storage.Transaction

	ctx, found := db.contexts[userAgent]
	if found {
		trn = ctx.tran
	}

	tokenizer := query.NewTokenizer(q)
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		return "", err
	}

	parser := query.NewParser(tokens)
	node, errs := parser.Parse()
	if len(errs) != 0 {
		// show first error message anyway...
		return "", errs[0]
	}

	analyzer := query.NewAnalyzer(db.catalog)
	analyzedQuery, err := analyzer.AnalyzeMain(node)
	if err != nil {
		return "", err
	}

	planner := query.NewPlanner(analyzedQuery)
	plan, _ := planner.PlanMain()

	executor := query.NewExecutor(db.storage, db.catalog, db.tranManager)
	result, err := executor.ExecuteMain(analyzedQuery, plan, trn)

	return db.stringfyResultSet(result), err
}

func (db *BogoDb) stringfyResultSet(r *meta.ResultSet) string{
	var sb strings.Builder

	sb.WriteString(r.Message)
	for _, c := range r.ColNames{
		sb.WriteString(c)
	}
	sb.WriteString("\n")

	for i, s := range r.Values{
		sb.WriteString(s)
		if i % len(r.ColNames) == 0{
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (db *BogoDb) Terminate() {
	if _, err := os.Stat(db.home); os.IsNotExist(err) {
		err := os.Mkdir(db.home, 0777)
		if err != nil {
			panic(err)
		}
	}

	err := storage.SaveCatalog(db.home, db.catalog)
	if err == nil {
		log.Printf("`catalog.db` has completely saved in %s\n", db.home)
	}

	err = db.storage.Terminate()
	if err == nil {
		log.Printf("data files have completely saved in %s\n", db.home)
	}

	os.Exit(0)
}
