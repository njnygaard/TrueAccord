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
const PaymentPlansRoute = "https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/payment_plans"
const PaymentsRoute = "https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/payments"

type DebtResponse struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
}

type PaymentPlanResponse struct {
	Id                   int     `json:"id"`
	DebtId               int     `json:"debt_id"`
	AmountToPay          float64 `json:"amount_to_pay"`
	InstallmentFrequency string  `json:"installment_frequency"`
	InstallmentAmount    float64 `json:"installment_amount"`
	StartDate            string  `json:"start_date"`
}

type PaymentResponse struct {
	PaymentPlanId int     `json:"payment_plan_id"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
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
	var payments []PaymentResponse
	var paymentPlans []PaymentPlanResponse

	debts, err = getDebts()
	if err != nil {
		logger.Fatal("error getting debts")
		return
	}

	payments, err = getPayments()
	if err != nil {
		logger.Fatal("error getting payments")
		return
	}

	paymentPlans, err = getPaymentPlans()
	if err != nil {
		logger.Fatal("error getting payment plans")
		return
	}

	logger.Info(spew.Sdump(debts))
	logger.Info(spew.Sdump(payments))
	logger.Info(spew.Sdump(paymentPlans))
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

func getPaymentPlans()(paymentPlans []PaymentPlanResponse, err error){

	logger := logrus.New()

	c := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, PaymentPlansRoute, nil)
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

	jsonErr := json.Unmarshal(body, &paymentPlans)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return
	}

	return

}


func getPayments()(payments []PaymentResponse, err error){

	logger := logrus.New()

	c := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, PaymentsRoute, nil)
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

	jsonErr := json.Unmarshal(body, &payments)
	if jsonErr != nil {
		logger.Error(jsonErr)
		return
	}

	return

}