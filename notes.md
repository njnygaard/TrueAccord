# Notes

[Instructions](https://gist.github.com/jeffling/2dd661ff8398726883cff09839dc316c)

> A description of your process and approach, including what you think you would have done differently given more time.

I started with scaffolding the utility with Cobra.
That really accelerates what I can do on the command line.

I didn't like how the API uses floats for currency.
I would ideally use unsigned integers for all amounts, signed if necessary.

I could combine my logic in the request layer to return interface decoders and decode when I have type information instead of making a typed get request function.
I have three functions that are very similar for getting the data from the API.
As that scales, I would definitely figure out a way to abstract that.

I didn't complete the testing because It would have taken a bit more time than I think I should have spent on this.
Specifically because I would have to think of test data for known edge cases.
Go makes it very easy to test but I didn't want to spend time designing tables.

> Some pointers on where to find the relevant logic in your code.

Everything relevant is in `TrueAccord/cmd/debts.go`.

The functions that interact with the API are prefaced with `get`.

`gatherResponses` is the high level aggregation and printing logic.

Functions prefaced with `process` transform data.

`calculateNextPaymentDueDate` is a helper function that necessarily does not error.

> Any design decisions or assumptions you made.

I expressly did not check for overpayment.

```go
// TODO: Negative Check
// If we decrement the payment plan below, it means the user overpaid.
// That opens up a whole other pathway, I think it is sufficient to say the plan is done for now.
```
