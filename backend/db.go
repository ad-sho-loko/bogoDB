package backend

import (
	"bogoDB/backend/query"
	"bogoDB/backend/storage"
	"log"
	"os"
	"os/signal"
)

type BogoDb struct {
	exit chan int
	storage *storage.Storage
	catalog *storage.Catalog
}

func NewBogoDb() (*BogoDb, error){
	// load the catalog if exists
	path, ok := os.LookupEnv("BOGO_HOME")
	if !ok{
		// for test codes
		path = "./tmp/"
	}

	catalog, err := storage.LoadCatalog(path)
	if err != nil{
		return nil, err
	}

	return &BogoDb{
		catalog:catalog,
		storage:storage.NewStorage(),
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
	parser := query.NewParser()
	node, err := parser.ParseMain(q)
	if err != nil{
		return err
	}

	analyzer := query.NewAnalyzer(db.catalog)
	analyzedQuery, err := analyzer.AnalyzeMain(node)
	if err != nil{
		return err
	}

	executor := query.NewExecutor(db.storage, db.catalog)
	return executor.ExecuteMain(analyzedQuery)
}

func (db *BogoDb) Terminate(){
	err := storage.SaveCatalog("./", db.catalog)
	if err == nil{
		log.Printf("catalog.db has completely saved in %s\n", "./")
	}
	os.Exit(0)
}
