package db

import (
	"github.com/ad-sho-loko/bogodb/query"
	"github.com/ad-sho-loko/bogodb/storage"
	"log"
	"os"
	"os/signal"
)

type BogoDb struct {
	exit chan int
	storage *storage.Storage
	catalog *storage.Catalog
	tranManager *storage.TransactionManager
}

func NewBogoDb() (*BogoDb, error){
	// load the catalog if exists
	path, ok := os.LookupEnv("BOGO_HOME")
	if !ok{
		// for test codes
		path = "."
	}

	catalog, err := storage.LoadCatalog(path)
	if err != nil{
		return nil, err
	}

	return &BogoDb{
		catalog:catalog,
		storage:storage.NewStorage(),
		tranManager:storage.NewTransactionManager(),
		exit:make(chan int, 1),
	}, nil
}

func (db *BogoDb) Init(){
	// set up the signal handler
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func(){
		for{
			r := <-sig
			if r == os.Interrupt{
				log.Println("ctrl+c detected, shutdown soon....")
				db.Terminate()
			}
		}
	}()
}

func (db *BogoDb) Execute(q string) error{
	tokenizer := query.NewTokenizer(q)
	tokens, err := tokenizer.Tokenize()
	if err != nil{
		return err
	}

	parser := query.NewParser(tokens)
	node, errs := parser.Parse()
	if len(errs) != 0{
		return errs[0] // show first error message anyway...
	}

	analyzer := query.NewAnalyzer(db.catalog)
	analyzedQuery, err := analyzer.AnalyzeMain(node)
	if err != nil{
		return err
	}

	planner := query.NewPlanner(analyzedQuery)
	plan, _ := planner.PlanMain()

	executor := query.NewExecutor(db.storage, db.catalog, db.tranManager)
	return executor.ExecuteMain(analyzedQuery, plan)
}

func (db *BogoDb) Terminate(){
	err := storage.SaveCatalog("./", db.catalog)
	if err == nil{
		log.Printf("catalog.db has completely saved in %s\n", "./")
	}
	os.Exit(0)
}
