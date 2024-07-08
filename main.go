package main

import (
	"database/sql"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// func resetInvoiceTotals(){

// 	}

func main() {

	app := pocketbase.New()
	///
	app.OnRecordBeforeCreateRequest("users", "invoices", "receipts", "receiptTypes").Add(
		func(e *core.RecordCreateEvent) error {

			type Number struct {
				Number int `db:"number" json:"number"`
			}

			var err error

			result := Number{}

			switch e.Collection.Name {
			case "users":
				err = app.Dao().DB().Select("number").From("users").OrderBy("number DESC").Limit(1).One(&result)
			case "receipts":
				err = app.Dao().DB().Select("number").From("receipts").OrderBy("number DESC").Limit(1).One(&result)
			case "invoices":
				err = app.Dao().DB().Select("number").From("invoices").OrderBy("number DESC").Limit(1).One(&result)
			case "receiptTypes":
				err = app.Dao().DB().Select("number").From("receiptTypes").OrderBy("number DESC").Limit(1).One(&result)
			default:
				return nil
			}

			if err != nil {
				switch err {
				case sql.ErrNoRows:
					e.Record.Set("number", 1)
				default:
					log.Fatal(err)
				}
				return nil
			}

			e.Record.Set("number", result.Number+1)

			return nil

		})

	app.OnRecordAfterCreateRequest("receipts").Add(func(e *core.RecordCreateEvent) error {

		value := e.Record.GetInt("value")
		invoiceId := e.Record.GetString("invoiceId")

		inv, err := app.Dao().FindRecordById("invoices", invoiceId)

		if err != nil {
			log.Panicln("error when adding receipt value to the invoice total value")
			log.Panicln("the error accur when reteiving the invoice")
			log.Panicln(err)

			return nil
		}

		inv.Set("totalValue", inv.GetInt("totalValue")+value)
		inv.Set("totalReceipts", inv.GetInt("totalReceipts")+1)
		inv.MarkAsNotNew()

		saveErr := app.Dao().Save(inv)
		if saveErr != nil {
			log.Panicln("error when adding receipt value to the invoice total value")
			log.Panicln("the error accur when saving the updated invoice")
			log.Panicln(saveErr)
			return nil
		}

		return nil
	})

	app.OnRecordAfterUpdateRequest("receipts").Add(func(e *core.RecordUpdateEvent) error {

		switch e.Record.GetBool("canceled") {
		case true:
			value := e.Record.GetInt("value")
			invId := e.Record.GetString("invoiceId")

			inv, err := app.Dao().FindRecordById("invoices", invId)

			if err != nil {
				log.Panicln("error when cancelling receipt value form the invoice total value")
				log.Panicln("when canceling receipt")
				log.Panicln("the error accur when reteiving the invoice")
				log.Panicln(err)

				return nil

			}
			inv.Set("totalValue", inv.GetInt("totalValue")-value)
			inv.Set("totalReceipts", inv.GetInt("totalReceipts")-1)
			inv.MarkAsNotNew()
			saveErr := app.Dao().Save(inv)

			if saveErr != nil {
				log.Panicln("error when removing receipt value form the invoice total value")
				log.Panicln("when canceling receipt")
				log.Panicln("the error accur when saving the updated invoice")
				log.Panicln(saveErr)

				return nil
			}
			return nil
		default:

			value := e.Record.GetInt("value")
			invId := e.Record.GetString("invoiceId")

			inv, err := app.Dao().FindRecordById("invoices", invId)

			if err != nil {
				log.Panicln("error when uncancelling receipt value form the invoice total value")
				log.Panicln("when un canceling receipt")
				log.Panicln("the error accur when reteiving the invoice")
				log.Panicln(err)

				return nil

			}
			inv.Set("totalValue", inv.GetInt("totalValue")+value)
			inv.Set("totalReceipts", inv.GetInt("totalReceipts")+1)
			inv.MarkAsNotNew()
			saveErr := app.Dao().Save(inv)

			if saveErr != nil {
				log.Panicln("error when inserting receipt value to  the invoice total value")
				log.Panicln("when un cancelling receipt")
				log.Panicln("the error accur when saving the updated invoice")
				log.Panicln(saveErr)

				return nil
			}
			return nil
		}

	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
