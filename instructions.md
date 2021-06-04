# TrueAccord Take Home

Please complete the following exercise at your leisure. There is no strict time limit, but if you've commited to complete this exercise at a certain time, please make your best attempt to do so, as our evaluators will be expecting it. 

To make it easier to evaluate the exercise, please complete the exercise in one of the following languages:

- Scala
- Java
- Go
- Typescript/Javascript
- Python

We will be assessing the code based off **correctness, readability, and extendability**. This exercise typically takes one to two hours with no distractions. Try to keep things simple.

## Problem Brief

At TrueAccord, we deal with customer debts in all of our products. The following problem will give you a database schema and API that outputs debts and payment plans. We expect you to think about the problem definitions, and try to solve them to the best of your ability.

## Problems

Imagine that you are writing an internal admin tool. Please complete the following tasks. Use the database schema and API documentation as a reference.

1.  Create a script that will consume data from the HTTP API endpoints described below and output debts to **stdout** in [JSON Lines format](https://jsonlines.org/). (NOTE: You are **not** expected to create your own server backend. Although the data is mocked, use the provided endpoints as though they serve real data.)
    - Each line should contain:
    - All the Debt object's fields returned by the API
    - An additional boolean value, "*is_in_payment_plan*", which is: 
      - true when the debt is associated with an active payment plan. 
      - false when there is no payment plan, or the payment plan is completed.
2. Provide a test suite that validates the output being produced, along with any other operations performed internally.
    - This can be done using any testing technique, but it should provide reasonable coverage of functionality.
3. Add a new field to the Debt objects in the output: "*remaining_amount*", containing the calculated amount remaining to be paid on the debt. Output the value as a JSON Number.
    - If the debt is associated with a payment plan, subtract from the payment plan's *amount_to_pay* instead. In exchange for signing up for a payment plan, we will allow them to pay a reduced amount to satisfy the debt.
    - All payments, whether on-time or not, contribute toward paying off a debt. 
4. Add a new field to the Debt object output: "*next_payment_due_date*", containing the ISO 8601 UTC **date** (i.e. “2019-09-07”) of when the next payment is due or null if there is no payment plan or if the debt has been paid off.
    - The *next_payment_due_date* can be calculated by using the payment plan *start_date* and *installment_frequency*. It should be the next installment date after the latest payment, even if this date is in the past.
    - The *next_payment_due_date* should be null if there is no payment plan or if the debt has been paid off.
    - Payments made on days outside the expected payment schedule still go toward paying off the *remaining_amount*, but do not change/delay the payment schedule.
## Submitting Your Work

Please email us the following once you are done:

- A .zip or .tar.gz of your solution
- A README.md file with the following information:
  - **Clear, step-by-step instructions** on how to run your application.
  - A high level overview of how you spent your time.
    - If you submit a git repo with numerous commits and clear commit messages, you may omit this.
  - A description of your process and approach, including what you think you would have done differently given more time.
  - Some pointers on where to find the relevant logic in your code. 
  - Any design decisions or assumptions you made. 
- Optionally, create an offline git repository and submit that once you're done.
    - By committing your work regularly with good comments, you'll have a high level overview of your work for free. This is not mandatory, just do what's easiest for you.

We will have two engineers review it independently after you submit it. 

## Simple HTTP API Service

This section details the three kinds of data objects and HTTP Service endpoints that provide access to them. Access each service endpoint using an **HTTP GET request**.

### Debts

A debt, which is money is owed to a collector.

* id (integer)
* amount (real) - amount owed in USD

#### HTTP Service URL: https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/debts

#### Example HTTP JSON Response

```json
[{
  "id": 0,
  "amount": 123.46
}]
```

### Payment Plans

A payment plan, which is an amount needed to resolve a debt, as well as the frequency of when it will be paid. Payment plans are associated with exactly one debt, and debts may not be associated with more than one payment plan.


* id (integer)
* debt_id (integer) - The associated debt.
* amount_to_pay (real) - Total amount (in USD) needed to be paid to resolve this payment plan. 
* installment_frequency (text) - The frequency of payments. Is one of: WEEKLY or BI_WEEKLY (14 days).
* installment_amount (real) - The amount (in USD) of each payment installment.
* start_date (string) - ISO 8601 date of when the first payment is due.

#### HTTP Service URL: https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/payment_plans

#### Example HTTP JSON Response Payload

```json
[{
  "id": 0,
  "debt_id": 0,
  "amount_to_pay": 102.50,
  "installment_frequency": "WEEKLY", 
  "installment_amount": 51.25,
  "start_date": "2020-09-28"
}]
```

### Payments

An individual payment installment which is made on a payment plan. Many-to-one with debts.

* payment_plan_id (integer)
* amount (real)
* date (string) - ISO 8601 date of when this payment occurred.

#### HTTP Service URL: https://my-json-server.typicode.com/druska/trueaccord-mock-payments-api/payments

#### Example HTTP JSON Response Payload

```json
[{
  "payment_plan_id": 0,
  "amount": 51.25,
  "date": "2020-09-29"
}]
```
