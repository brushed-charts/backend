package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brushed-charts/backend/lib/cloudlogging"
)

const (
	serviceName    = "institution/oanda"
	projectID      = "brushed-charts"
	envTokenName   = "OANDA_API_TOKEN"
	minRefreshRate = "1s"
)

var (
	apiURL           string
	bqPriceShortterm string
	bqPriceArchive   string
	bigQueryDataset  string
	watchlistPath    string
	latestCandlePath string
)

func main() {
	apiURL = os.Getenv("OANDA_API_URL")
	bqPriceShortterm = os.Getenv("OANDA_BIGQUERY_SHORTTERM_TABLENAME")
	bqPriceArchive = os.Getenv("OANDA_BIGQUERY_ARCHIVE_TABLENAME")
	bigQueryDataset = os.Getenv("OANDA_BIGQUERY_DATASET")
	watchlistPath = os.Getenv("OANDA_WATCHLIST_PATH")
	latestCandlePath = os.Getenv("OANDA_LATEST_CANDLE_PATH")

	setAPIKeyEnvVariable()

	err := checkEnvVariable()
	if err != nil {
		cloudlogging.ReportCritical(cloudlogging.EntryFromError(err))
		log.Fatalf("%v", err)
	}

	// bigqueryDeduplication()

	id, err := getAccountID()
	if err != nil {
		log.Fatal("Can't retrieve the accound ID")
	}

	instruments := parseWatchlist()
	go watchlistRefreshLoop(minRefreshRate, &instruments)

	stream, err := fetchlatestCandles(id, &instruments, minRefreshRate, []string{"S5", "M1", "H1", "D"})
	if err != nil {
		log.Fatal(err)
	}

	go streamToBigQuery(stream.candles)

	for {
		select {
		case err := <-stream.err:
			log.Printf("%v", err)
		case err := <-stream.fatal:
			mess := fmt.Errorf("Fatal Error program exited : %v", err)
			cloudlogging.ReportEmergency(cloudlogging.EntryFromError(mess))
			log.Fatalf("%v", err)
		}
	}

}

func setAPIKeyEnvVariable() {
	// Check for OANDA API KEY
	if _, b := os.LookupEnv(envTokenName); !b {
		token, err := getToken()
		if err != nil {
			log.Fatal(err)
		}

		// Set OANDA API KEY in env variable
		os.Setenv(envTokenName, token)
	}
}

func checkEnvVariable() error {
	baseErr := "Environnement variables are empty : "
	errSentence := baseErr
	if apiURL == "" {
		errSentence += "\napiURL"
	}

	if bqPriceArchive == "" {
		errSentence += "\nbqTableNameArchive"
	}

	if bqPriceShortterm == "" {
		errSentence += "\nbqTableNameShortterm"
	}

	if watchlistPath == "" {
		errSentence += "\nwatchlistPath"
	}

	if latestCandlePath == "" {
		errSentence += "\bLatestCandlePath"
	}

	if os.Getenv(envTokenName) == "" {
		errSentence += "\n" + envTokenName
	}

	if errSentence != baseErr {
		err := fmt.Errorf(errSentence)
		return err
	}

	return nil
}
