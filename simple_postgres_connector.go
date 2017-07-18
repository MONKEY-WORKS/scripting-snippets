/*
 *  Simple PostgreSQL connector to read data provided in a data base.
 *  This connector just reads the last line of a table
 *  Author: Johannes Tandler
 *  Date: 18.07.2017
 *
 *
 */
package scripts

import (
	// import for simple print lines and so on
	"fmt"

	// import of base sql framework
	"database/sql"

	// import for postgresql driver
	// We import the package solely for its side-effects (initialization)
	// For this we have an underscore at the beginning of the line
	_ "github.com/lib/pq"

	// import of application related stuff
	"git.monkey-works.de/scripting/api"
	"monkey-works.de/model"

	// import for time related things
	"time"
)

var (
	// global database reference
	db *sql.DB

	// global reference to application content
	app model.Application
)

// @script
func InitializeScripting(application model.Application) {
	fmt.Println("Hello Scripting")
	app = application

	// get status item
	status, ok := application.ClientDataModel().FindDataItemByName("status").(model.StringDataItem)
	if !ok {
		panic("DataItem \"status\" not found")
	}

	// establish connection to postgressql database
	var err error
	db, err = sql.Open("postgres", "host=192.168.1.66 user=postgres dbname=postgres password=mysecretpassword sslmode=disable")

	if err != nil {
		// error occured, notify user in status text
		status.SetCurrentValue(err.Error())
		return
	}

	//refresh values
	go refreshValues()

	// listen to refresh button and trigger refresh if button get pressed
	registerButton()
}

// refreshValues refresh all DataItems by retrieving last entry of a postgressql library
func refreshValues() {
	// get status data item in order to propagate error messages
	status := app.ClientDataModel().FindDataItemByName("status").(model.StringDataItem)

	// query database
	rows, err := db.Query("SELECT id, value, timestamp FROM \"sample-data\" ORDER BY id DESC LIMIT 1")
	if err != nil {
		// catch errors
		status.SetCurrentValue(err.Error())
		return
	}
	// close database connection after this method
	defer rows.Close()

	// get first row
	rows.Next()

	var id int32
	var value float32
	var timestamp time.Time

	// parse data
	if err = rows.Scan(&id, &value, &timestamp); err != nil {
		status.SetCurrentValue(err.Error())
		return
	}

	// get data items
	idDataItem := app.ClientDataModel().FindDataItemByName("id").(model.IntegerDataItem)
	valueDataItem := app.ClientDataModel().FindDataItemByName("value").(model.NumberDataItem)
	timestampDataItem := app.ClientDataModel().FindDataItemByName("timestamp").(model.StringDataItem)

	// set new data item values
	idDataItem.SetCurrentValue(id)
	valueDataItem.SetCurrentValue(float64(value))
	timestampDataItem.SetCurrentValue(timestamp.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	status.SetCurrentValue("Retrieved data successfully")
}

// registerButton registers a listener to the refresh button, so we can refresh all values when the button gets triggered
func registerButton() {
	// get button reference
	btn := app.ClientDataModel().FindDataItemByName("refreshTriggered").(model.BooleanDataItem)
	btn.SetCurrentValue(false)

	// register adapter
	btn.AddAdapter(func(change api.FeatureChange) {
		val := btn.CurrentValue()

		// if button was pressed
		if val {
			// reset button state
			btn.SetCurrentValue(false)

			// refresh value
			refreshValues()
		}
	})
}
