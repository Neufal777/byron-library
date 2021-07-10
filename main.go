package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Byron/backend"
	"github.com/Byron/data"
	"github.com/Byron/executions"
	"github.com/Byron/parsecore"
	"github.com/gorilla/mux"
	"github.com/ttacon/chalk"
)

func init() {
	os.Setenv("MONGO_CONNECTION", "mongodb+srv://byron:Black_nebula000@byron.dqvrs.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	os.Setenv("MONGO_DATABASE", "byron")
	os.Setenv("MONGO_COLLECTION", "byronArticles")

	/*
		Setting variables
	*/
	log.Println(chalk.Magenta.Color("Setting up variables"))
	log.Println(chalk.Green.Color("Mongo conenction: " + os.Getenv("MONGO_CONNECTION")))
	log.Println(chalk.Green.Color("Mongo database: " + os.Getenv("MONGO_DATABASE")))
	log.Println(chalk.Green.Color("Mongo collection: " + os.Getenv("MONGO_COLLECTION")))

}

func main() {
	execMode := flag.String("exec", "web", "Select mode of execution")
	flag.Parse()

	switch *execMode {
	case "web":
		WebExecute()
	case "parse":
		ParseExecute()
	case "parseLibgen":
		ParseExecuteLibgen()
	case "insert":
		data.InsertArticlesDatabase()
	case "delete":
		data.DeleteDuplicates("Inventory/")
	case "test":
		TestingExecute()
	case "proxy":
		ProxyTesting()
	default:
		log.Println("Not found:", *execMode)
		WebExecute()
	}
}

func ProxyTesting() {
	parsecore.ProxyScraping("https://es.wikipedia.org/wiki/Los_%C3%81ngeles")
}
func TestingExecute() {
	// for i := 0; i < 200; i++ {
	// 	fmt.Println("go LibrospdfgratismundoExecution(" + strconv.Itoa(i*5) + ", " + strconv.Itoa((i*5)+5) + ")")
	// }

	//executions.ArchiveOrgExecution()
	//data.RecoverSource("Inventory/", "libgen.is")

}
func ParseExecuteLibgen() {
	go executions.LibGenExecution()
	go executions.LibGenExecution01()
	go executions.LibGenExecution02()
	go executions.LibGenExecution03()
	go executions.LibGenExecution04()
	go executions.LibGenExecution05()
	go executions.LibGenExecution06()
	go executions.LibGenExecution07()
	go executions.LibGenExecution08()
	go executions.LibGenExecution09()
	go executions.LibGenExecution10()
	go executions.LibGenExecution11()
	go executions.LibGenExecution12()
	go executions.LibGenExecution13()
	go executions.LibGenExecution13()
	go executions.LibGenExecution14()
	executions.LibGenExecution15()
}

func ParseExecute() {
	go executions.BookRixExecution()
	go executions.FreeditorialExecution()
	go executions.ManyBooksExecution()
	go executions.LibrospdfgratismundoExecution()
	go executions.OpenlibraExecution()

	executions.ArchiveOrgExecution()

}

func WebExecute() {
	r := mux.NewRouter()

	//loading files from assets folder
	fh := http.FileServer(http.Dir("./web/assets/"))
	r.PathPrefix("/web/assets/").Handler(http.StripPrefix("/web/assets/", fh))

	//web handlers
	r.HandleFunc("/search", backend.SearchResults).Methods("GET")
	r.HandleFunc("/", backend.HomeSearchBar)

	//identify and assign PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	log.Println("PORT=> ", port)
	http.ListenAndServe(":"+port, r)
}
