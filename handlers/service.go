package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type HowMuchResponse struct {
	NextPayDay      time.Time
	DaysUntilPayDay int
}

func ValidateEndPointURL(url string) bool {
	// Validate the url endpoints
	// It returns true if the url is ok, else it returns false

	//verify first endpoint url (how-much)
	tillSalaryHowMuchRegex, err := regexp.Compile(`/till-salary/how-much`)
	if err != nil {
		fmt.Printf("Error while validating endpoint url: %v", err)
		return false
	}

	//verify second endpoint url (pay-day)
	listPayDayDatesRegex, err := regexp.Compile(`^/till-salary/pay-day/([1-9]|[1-3][0-9])/list-dates$`)
	if err != nil {
		fmt.Printf("Error while validating endpoint url: %v", err)
		return false
	}

	if tillSalaryHowMuchRegex.MatchString(url) || listPayDayDatesRegex.MatchString(url) {
		return true
	}

	return false
}

func CheckIfMonthHas31Days(payDay int, month time.Month) time.Time {
	// Check if month from date sent has 31 days or not
	// If it has not I take the closest date that is not in the weekend
	// It returns the correct date

	currYear := time.Now().Year()

	if payDay == 31 {
		//I take the last day of the current month to see if it is 31 or not
		auxTime := time.Date(currYear, month+1, 0, 0, 0, 0, 0, time.Local)

		//check if is in weekend or not, because in weekend you can not be paid
		if auxTime.Day() != 31 {
			if auxTime.Weekday() == time.Saturday {
				payDay = auxTime.Day() - 1
			} else if auxTime.Weekday() == time.Sunday {
				payDay = 1
				month++
			} else {
				payDay = auxTime.Day()
			}
		}
	}

	nextPayDay := time.Date(currYear, month, payDay, 0, 0, 0, 0, time.Local)

	return nextPayDay
}

func CalculateDaysUntilPayday(payDay int) HowMuchResponse {
	// Calculate the next pay day date and the number of days until
	// It returns a struct with two fields: number of days until pay day and the pay day date
	// I made a struct to be easier to get either the number of days or the date

	now := time.Now()
	month := now.Month()
	var howMuchResponse HowMuchResponse

	howMuchResponse.NextPayDay = CheckIfMonthHas31Days(payDay, month)

	if howMuchResponse.NextPayDay.Before(now) {
		howMuchResponse.NextPayDay = time.Date(now.Year(), now.Month()+1, payDay, 0, 0, 0, 0, time.Local)
	}

	howMuchResponse.DaysUntilPayDay = int(howMuchResponse.NextPayDay.Sub(now).Hours()/24) + 1

	return howMuchResponse
}

func CalculatePayDayDates(payDay int) []string {
	// Calculate the pay day dates from this year
	// It returns a string slice with all the pay day dates

	now := time.Now()

	var dates []string
	for month := now.Month(); month <= time.December; month++ {
		checkedDate := CheckIfMonthHas31Days(payDay, month)
		if checkedDate.After(now) {
			dates = append(dates, checkedDate.Format("2006-01-02"))
		}
	}

	return dates
}

func TillSalaryHandler(writer http.ResponseWriter, req *http.Request) {
	// Handler function for the next pay day

	//check if is a Get method
	if req.Method != http.MethodGet {
		fmt.Printf("Method not allowed")
		writer.WriteHeader(http.StatusMethodNotAllowed)

		err := json.NewEncoder(writer).Encode(Response{Message: "Method not allowed"})
		if err != nil {
			fmt.Printf("Error while writing response message for method not allowed case: %v", err)
			return
		}

		return
	}

	//validate url
	if !ValidateEndPointURL(req.URL.Path) {
		fmt.Printf("Error while validating endpoint url")
		writer.WriteHeader(http.StatusBadRequest)

		err := json.NewEncoder(writer).Encode(Response{Message: "Invalid how-much url"})
		if err != nil {
			fmt.Printf("Error while writing response message for invalid how-much url case: %v", err)
			return
		}

		return
	}

	//get the pay day number parameter and check if is ok
	payDayStr := req.URL.Query().Get("pay_day")
	payDay, err := strconv.Atoi(payDayStr)
	if err != nil || payDay < 1 || payDay > 31 {
		writer.WriteHeader(http.StatusBadRequest)

		err = json.NewEncoder(writer).Encode(Response{Message: "Invalid pay_day parameter"})
		if err != nil {
			fmt.Printf("Error while writing response message for invalid pay_day parameter case: %v", err)
			return
		}

		return
	}

	//calculate the number of days and the date
	howMuchResponse := CalculateDaysUntilPayday(payDay)

	//write next pay day, number of days until pay day and the message
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(Response{Data: map[string]interface{}{
		"next_pay_day":       howMuchResponse.NextPayDay.Format("2006-01-02"),
		"days_until_pay_day": howMuchResponse.DaysUntilPayDay,
	},
		Message: "Days until pay day"})
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(Response{Message: "Error while writing data"})
		if err != nil {
			fmt.Printf("Error while writing response message: %v", err)
			return
		}
		return
	}
}

func ListPayDayDatesHandler(writer http.ResponseWriter, req *http.Request) {
	// Handler function for the pay day dates

	//check if is a Get method
	if req.Method != http.MethodGet {
		fmt.Printf("Method not allowed")
		writer.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(writer).Encode(Response{Message: "Method not allowed"})
		if err != nil {
			fmt.Printf("Error while writing response message for method not allowed case: %v", err)
			return
		}
		return
	}

	//validate url
	if !ValidateEndPointURL(req.URL.Path) {
		fmt.Printf("Error while validating endpoint url")
		writer.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(writer).Encode(Response{Message: "Invalid pay-day url"})
		if err != nil {
			fmt.Printf("Error while writing response message for invalid pay-day url case: %v", err)
			return
		}
		return
	}

	//get the pay day number
	var payDayStr string
	position := len("/till-salary/pay-day/")
	if req.URL.Path[position+1] == '/' {
		payDayStr = string([]byte{req.URL.Path[position]})
		position++
	} else {
		payDayStr = req.URL.Path[position : position+2]
		position = position + 2
	}

	//check if the parameter pay day is ok (it exists, it is a number between 1 and 31)
	payDay, err := strconv.Atoi(payDayStr)
	if err != nil || payDay < 1 || payDay > 31 {
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(Response{Message: "Invalid pay_day parameter"})
		if err != nil {
			fmt.Printf("Error while writing response message for invalid pay_day parameter case: %v", err)
			return
		}
		return
	}

	//calculate the dates
	dates := CalculatePayDayDates(payDay)

	//write next pay day dates from this year
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(Response{Data: map[string]interface{}{
		"pay_day_dates": dates,
	},
		Message: "Pay day dates"})
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(Response{Message: "Error while writing data"})
		if err != nil {
			fmt.Printf("Error while writing data case: %v", err)
			return
		}
		return
	}
}
