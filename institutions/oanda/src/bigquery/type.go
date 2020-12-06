package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
)

// CandleOHLC ss
type CandleOHLC struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// Save ss
func (c *CandleOHLC) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"open":  c.Open,
		"high":  c.High,
		"low":   c.Low,
		"close": c.Close,
	}, "", nil
}

// CandleRow ssd
type CandleRow struct {
	Instrument  string     `json:"instrument"`
	Date        string     `json:"date"`
	Granularity string     `json:"granularity"`
	Bid         CandleOHLC `json:"bid"`
	Ask         CandleOHLC `json:"ask"`
	Volume      int        `json:"volume"`
}

// Save sdsd
func (c *CandleRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"instrument":  c.Instrument,
		"date":        c.Date,
		"granularity": c.Granularity,
		"bid":         c.Bid,
		"ask":         c.Ask,
		"volume":      c.Volume,
	}, "", nil
}

// CandleLink concat information needed to insert candles
// with bigquery. You need to call init() after initializing
// the main fields
type CandleLink struct {
	client              *bigquery.Client
	context             context.Context
	inserters           []*bigquery.Inserter
	Dataset             string         `json:"dataset"`
	TablePriceArchive   string         `json:"tablePriceArchive"`
	TablePriceShortterm string         `json:"tablePriceShortterm"`
	GCPProjectID        string         `json:"gcpProjectID"`
	InputCandleStream   chan CandleRow // InputCandleStream is already initialized when calling init()
	OutputError         chan error     // OutputError is already initialized when calling init()
}

// Init initialize InputCandleStream and OutputError channels
// but also it initialize some hidden fields (client, context, inserters)
func (c *CandleLink) Init() error {
	err := c.initClient()
	if err != nil {
		return err
	}
	c.InputCandleStream = make(chan CandleRow)
	c.OutputError = make(chan error)
	c.inserters = append(c.inserters, c.client.Dataset(c.Dataset).Table(c.TablePriceArchive).Inserter())
	c.inserters = append(c.inserters, c.client.Dataset(c.Dataset).Table(c.TablePriceShortterm).Inserter())
	return nil
}

func (c *CandleLink) initClient() error {
	var err error
	c.context = context.Background()
	c.client, err = bigquery.NewClient(c.context, c.GCPProjectID)
	if err != nil {
		err = errors.New(fmt.Sprintf("bigquery.NewClient: %v", err))
		return err
	}
	return nil
}

func (c *CandleRow) equalIgnoreDate(row CandleRow) bool {
	if c.Instrument != row.Instrument {
		return false
	}
	if c.Granularity != row.Granularity {
		return false
	}
	return true
}