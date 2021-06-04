package cmd

import (
	"encoding/json"
	"fmt"
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
	Id                 int       `json:"id"`
	Amount             float64   `json:"amount"`
	IsInPaymentPlan    bool      `json:"is_in_payment_plan"`
	RemainingAmount    float64   `json:"remaining_amount"`
	NextPaymentDueDate time.Time `json:"next_payment_due_date"`
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

	var processedDebts []Debt
	processedDebts, err = processDebts(debts, processedPaymentPlans, payments)
	if err != nil {
		logger.Error("could not process debts")
		return
	}

	for i := range processedDebts {

		j, err := json.Marshal(processedDebts[i])
		if err != nil {
			logger.Error("Failed to generate json", err)
		}

		// Bare
		fmt.Printf("%s\n", string(j))

		// Decorated
		//logger.Info(fmt.Sprintf("%s\n", string(j)))
	}

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

func processDebts(d []DebtResponse, pp []PaymentPlan, _ []PaymentResponse) (debts []Debt, err error) {

	// DONE: Consume PaymentPlan not PaymentPlanResponse
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
		var remainingAmount float64
		var nexPaymentDueDate time.Time

		for j := range pp {

			// If there is a payment plan, we use the amount_to_pay instead of the debt amount
			// There is a discount for signing up for a payment plan
			if d[i].Id == pp[j].DebtId {
				paymentPlanFound = true

				if pp[j].IsComplete == true {
					remainingAmount = 0
				} else {
					remainingAmount = pp[j].AmountToPay
					// "Payments made on days outside the expected payment schedule still go toward paying off the *remaining_amount*, but do not change/delay the payment schedule."
					// I take that to mean that late or early payments do not change next_payment_due_date
					// So I will not factor in any of the dates on payment objects.
					// That breaks down the calculation in to a modulus using an enum for the frequency.
					nexPaymentDueDate = calculateNextPaymentDueDate(pp[j].StartDate, pp[j].InstallmentFrequency)
				}

				break
			}
		}

		var debt Debt

		debt.Id = d[i].Id
		debt.Amount = d[i].Amount
		debt.IsInPaymentPlan = paymentPlanFound
		debt.NextPaymentDueDate = nexPaymentDueDate

		if paymentPlanFound {
			debt.RemainingAmount = remainingAmount
		} else {
			debt.RemainingAmount = d[i].Amount
		}

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
		// I didn't complete this because I think it is unnecessary.
		// If we decrement the payment plan below, it means the user overpaid.
		// That opens up a whole other pathway, I think it is sufficient to say the plan is done for now.
		if pp[i].AmountToPay <= 0 {
			paymentPlan.IsComplete = true
		}

		paymentPlans = append(paymentPlans, paymentPlan)
	}

	return

}

func calculateNextPaymentDueDate(startDate string, installmentFrequency string) (nextPaymentDueDate time.Time) {

	const (
		format   = "2006-01-02"
		weekly   = "WEEKLY"
		biweekly = "BI_WEEKLY"
	)

	var err error
	var start time.Time

	now := time.Now()

	start, err = time.Parse(format, startDate)
	if err != nil {
		return time.Time{}
	}

	var freq time.Duration
	switch installmentFrequency {
	case weekly:
		freq = time.Hour * 24 * 7
		break
	case biweekly:
		freq = time.Hour * 24 * 7 * 2
		break
	default:
		return time.Time{}
	}

	var nextFound bool

	for !nextFound {
		start = start.Add(freq)
		if start.After(now) {
			nextFound = true
			nextPaymentDueDate = start
		}
	}

	return
}