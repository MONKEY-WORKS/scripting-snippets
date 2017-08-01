package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)


//The endpoint of our web service we want to connect to
//Change this to your endpoint
const URL string = "http://localhost:8080/SPSSimulator/rest/spssimulator"
const USERNAME string = "user"
const PASSWORD string = "password"
const DATABASE string = "database"
const PORT string = "1234"
const QUERY string = "Select name, value, timestamp FROM \"data_table\""

//The interval at which the connector is updated
//Time is defined in milliseconds
const UPDATE_INTERVAL time.Duration = 1000

// Global channel to consume received data
var Data = make(chan *Response)

// Global channel to control periodic task
var Stop = make(chan bool)


var db *sql.DB
// Data structure of simple data item provided by REST interface
// Adapt this to your data structure
type DataItem struct {
	Name      string  `json:name`
	Value     float32 `json:value`
	Timestamp string  `json:timestamp`
}

// The repsonse of REST interface
// Adapt this to your data structure
type Response struct {
	DataItems []DataItem `json:"dataItems`
}

// starts a thread to receive data periodically and consumes messages from the global channels
func main() {
	go getData()
	fmt.Println("Subscription started")
	for {
		select {
		case val := <-Data:
			processData(val)
		case <-Stop:
			fmt.Println("Stop")
			return
		}
	}
}

// processing of received data
func processData(data *Response) {
	if data == nil {
		fmt.Println("Error")
		return
	}
	for _, value := range data.DataItems {
		fmt.Println("DataItem ", value.Name, " Value ", value.Value, " Timestamp ", value.Timestamp)
		//
		// Add your code here ...
		//
	}
}



// starts every UPDATE_INTERVAL milliseconds a  Request to get data from the data base
// the response is forwarded to the global channel if available
func getData() {
	var err error
	db, err = sql.Open("postgres", "host="+URL+":"+PORT+"user="+USERNAME+" dbname="+DATABASE+" password="+PASSWORD+" sslmode=disable")
	if(err!=nil) {
		fmt.Println("Connection error")
		Stop <- true
		return
	}
	var resp Response
	for range time.Tick(time.Millisecond * UPDATE_INTERVAL) {
		rows, err := db.Query(QUERY)
		if(err != nil) {
			//some error happened
			Data <- nil
			continue
		}


		for rows.Next {
			var dataItem DataItem
			rows.Scan(&dataItem.Name,&dataItem.Value,&dataItem.Timestamp) 
			//	Data <- nil
			//	continue
			//}
			resp.DataItems = append(resp.DataItems,dataItem)

		}
		Data <- &resp

	}
	Stop <- true
}