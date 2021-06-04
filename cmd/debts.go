package cmd

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const DebtsRoute = "https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/debts"

type DebtResponse struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
}

var debtsCmd = &cobra.Command{
	Use:   "debts",
	Short: "output debts to stdout in JSON Lines format",
	Long: `
Consume data from the HTTP API endpoints described below and output debts to stdout in JSON Lines format.
- Each line contains:
    - All the Debt object's fields returned by the API
    - An additional boolean value, "is_in_payment_plan", which is: 
      - true when the debt is associated with an active payment plan. 
      - false when there is no payment plan, or the payment plan is completed.`,
	Run: func(cmd *cobra.Command, args []string) {
		gatherResponses()
	},
}

func init() {
	rootCmd.AddCommand(debtsCmd)
}

func gatherResponses()(){

	logger := logrus.New()
	var err error
	var debts []DebtResponse

	debts, err = getDebts()
	if err != nil {
		logger.Fatal("error getting debts")
		return
	}

	logger.Info(spew.Sdump(debts))
}

func getDebts()(debts []DebtResponse, err error){

	logger := logrus.New()

	c := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, DebtsRoute, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	res, getErr := c.Do(req)
	if getErr != nil {
		logger.Error(getErr)
		return
	}

	if res.Body != nil {
		// Closure to explicitly ignore deferred error
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logger.Error(readErr)
		return
	}

	jsonErr := json.Unmarshal(body, &debts)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return
	}

	return

}