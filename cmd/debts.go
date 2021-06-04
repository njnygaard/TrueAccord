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

type Debt struct {
	Id              int     `json:"id"`
	Amount          float64 `json:"amount"`
	IsInPaymentPlan bool    `json:"is_in_payment_plan"`
}

type PaymentPlanResponse struct {
	Id                   int     `json:"id"`
	DebtId               int     `json:"debt_id"`
	AmountToPay          float64 `json:"amount_to_pay"`
	InstallmentFrequency string  `json:"installment_frequency"`
	InstallmentAmount    float64 `json:"installment_amount"`
	StartDate            string  `json:"start_date"`
}

type PaymentPlan struct {
	Id                   int     `json:"id"`
	DebtId               int     `json:"debt_id"`
	AmountToPay          float64 `json:"amount_to_pay"`
	InstallmentFrequency string  `json:"installment_frequency"`
	InstallmentAmount    float64 `json:"installment_amount"`
	StartDate            string  `json:"start_date"`
	IsComplete           bool    `json:"is_complete"`
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

func gatherResponses() {

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

	var processedPaymentPlans []PaymentPlan
	processedPaymentPlans, err = processPaymentPlans(paymentPlans, payments)
	if err != nil {
		logger.Error("could not process debts")
		return
	}

	//logger.Info(spew.Sdump(processedPaymentPlans))

	var processedDebts []Debt
	processedDebts, err = processDebts(debts, processedPaymentPlans, payments)
	if err != nil {
		logger.Error("could not process debts")
		return
	}

	logger.Info(spew.Sdump(processedDebts))
}

func processDebts(d []DebtResponse, pp []PaymentPlan, _ []PaymentResponse) (debts []Debt, err error) {

	// TODO: Consume PaymentPlan not PaymentPlanResponse
	// We need to know if a given paymentPlan is complete

	/*
	   - An additional boolean value, "*is_in_payment_plan*", which is:
	     - true when the debt is associated with an active payment plan.
	     - false when there is no payment plan, or the payment plan is completed.
	*/

	// That's interesting. I should decorate the payment plan first with a field to determine
	// if it is complete. I have to get all the payments and add them up to see if a given
	// paymentPlan is complete.

	for i := range d {
		var paymentPlanFound bool // zero value false
		for j := range pp {
			if d[i].Id == pp[j].DebtId && pp[j].IsComplete != true {
				paymentPlanFound = true
				break
			}
		}

		var debt Debt
		debt.Id = d[i].Id
		debt.Amount = d[i].Amount

		debt.IsInPaymentPlan = paymentPlanFound

		debts = append(debts, debt)
	}

	return
}

func processPaymentPlans(pp []PaymentPlanResponse, p []PaymentResponse) (paymentPlans []PaymentPlan, err error) {

	for i := range p {
		for j := range pp {
			if p[i].PaymentPlanId == pp[j].Id {
				pp[j].AmountToPay -= p[i].Amount
			}
		}
	}

	for i := range pp {

		var paymentPlan PaymentPlan
		paymentPlan.Id = pp[i].Id
		paymentPlan.DebtId = pp[i].DebtId
		paymentPlan.AmountToPay = pp[i].AmountToPay
		paymentPlan.InstallmentFrequency = pp[i].InstallmentFrequency
		paymentPlan.InstallmentAmount = pp[i].InstallmentAmount
		paymentPlan.StartDate = pp[i].StartDate

		// TODO: Negative Check
		if pp[i].AmountToPay <= 0 {
			paymentPlan.IsComplete = true
		}

		paymentPlans = append(paymentPlans, paymentPlan)
	}

	return

}

func getDebts() (debts []DebtResponse, err error) {

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

func getPaymentPlans() (paymentPlans []PaymentPlanResponse, err error) {

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

func getPayments() (payments []PaymentResponse, err error) {

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
