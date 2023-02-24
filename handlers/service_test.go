package handlers

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateEndPointURL(t *testing.T) {
	type testCases struct {
		name string
		url  string
		want bool
	}
	for _, scenario := range []testCases{
		{
			name: "how-much url ok",
			url:  "/till-salary/how-much?pay_day=31",
			want: true,
		},
		{
			name: "pay-day url ok",
			url:  "/till-salary/pay-day/15/list-dates",
			want: true,
		},
		{
			name: "invalid url",
			url:  "/invalid-url",
			want: false,
		},
		{
			name: "empty url",
			url:  "",
			want: false,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			got := ValidateEndPointURL(scenario.url)

			if got != scenario.want {
				t.Errorf("wanted: %v, got: %v", scenario.want, got)
			}
		})
	}
}

func TestCheckIfMonthHas31Days(t *testing.T) {
	type testCases struct {
		name        string
		payDayInput int
		monthInput  time.Month
		want        time.Time
	}
	for _, scenario := range []testCases{
		{
			name:        "31 January",
			payDayInput: 31,
			monthInput:  time.January,
			want:        time.Date(time.Now().Year(), time.January, 31, 0, 0, 0, 0, time.Local),
		},
		{
			name:        "31 February",
			payDayInput: 31,
			monthInput:  time.February,
			want:        time.Date(time.Now().Year(), time.February, 28, 0, 0, 0, 0, time.Local),
		},
		{
			name:        "31 April",
			payDayInput: 31,
			monthInput:  time.April,
			want:        time.Date(time.Now().Year(), time.May, 1, 0, 0, 0, 0, time.Local),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			got := CheckIfMonthHas31Days(31, scenario.monthInput)
			if !got.Equal(scenario.want) {
				t.Errorf("wanted: %v, got: %v", scenario.want, got)
			}
		})
	}
}

// I verify only the date because the number of days until the pay day
// will be different according to the day you run
func TestCalculateDaysUntilPayday(t *testing.T) {
	type testCases struct {
		name             string
		payDayInput      int
		wantedNextPayDay time.Time
	}
	for _, scenario := range []testCases{
		{
			name:             "pay-day: 5",
			payDayInput:      5,
			wantedNextPayDay: time.Date(time.Now().Year(), time.March, 5, 0, 0, 0, 0, time.Local),
		},
		{
			name:             "pay-day: 10",
			payDayInput:      10,
			wantedNextPayDay: time.Date(time.Now().Year(), time.March, 10, 0, 0, 0, 0, time.Local),
		},
		{
			name:             "pay-day: 28 (this month)",
			payDayInput:      28,
			wantedNextPayDay: time.Date(time.Now().Year(), time.February, 28, 0, 0, 0, 0, time.Local),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			got := CalculateDaysUntilPayday(scenario.payDayInput)

			if !got.NextPayDay.Equal(scenario.wantedNextPayDay) {
				t.Errorf("Wanted: %v, got: %v", scenario.wantedNextPayDay, got.NextPayDay)
			}
		})
	}
}

func TestCalculatePayDayDates(t *testing.T) {
	type testCases struct {
		name        string
		payDayInput int
		want        []string
	}
	for _, scenario := range []testCases{
		{
			name:        "pay-day: 15",
			payDayInput: 15,
			want: []string{time.Date(time.Now().Year(), time.March, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.April, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.May, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.June, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.July, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.August, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.September, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.October, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.November, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.December, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02")},
		},
		{
			name:        "pay-day: 5",
			payDayInput: 5,
			want: []string{time.Date(time.Now().Year(), time.March, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.April, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.May, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.June, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.July, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.August, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.September, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.October, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.November, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				time.Date(time.Now().Year(), time.December, 5, 0, 0, 0, 0, time.Local).Format("2006-01-02")},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			got := CalculatePayDayDates(scenario.payDayInput)

			if cmp.Diff(scenario.want, got) != "" {
				t.Errorf("Wanted: %#v, got: %#v", scenario.want, got)
			}
		})
	}
}

// I use CalculateDaysUntilPayday function because it is already tested
// Here, first, I verify the status code and then the body
func TestTillSalaryHandler(t *testing.T) {
	type testCases struct {
		name         string
		methodType   string
		url          string
		statusInput  int
		payDay       int
		wantResponse Response
	}
	for _, scenario := range []testCases{
		{
			name:        "Valid URL (12)",
			methodType:  http.MethodGet,
			url:         "/till-salary/how-much?pay_day=12",
			statusInput: http.StatusCreated,
			payDay:      12,
			wantResponse: Response{Data: map[string]interface{}{
				"next_pay_day":       CalculateDaysUntilPayday(12).NextPayDay.Format("2006-01-02"),
				"days_until_pay_day": float64(CalculateDaysUntilPayday(12).DaysUntilPayDay),
			}, Message: "Days until pay day"},
		},
		{
			name:        "Valid URL (27)",
			methodType:  http.MethodGet,
			url:         "/till-salary/how-much?pay_day=27",
			statusInput: http.StatusCreated,
			payDay:      27,
			wantResponse: Response{Data: map[string]interface{}{
				"next_pay_day":       CalculateDaysUntilPayday(27).NextPayDay.Format("2006-01-02"),
				"days_until_pay_day": float64(CalculateDaysUntilPayday(27).DaysUntilPayDay),
			}, Message: "Days until pay day"},
		},
		{
			name:         "Invalid URL",
			methodType:   http.MethodGet,
			url:          "/invalid-url?pay_day=12",
			statusInput:  http.StatusBadRequest,
			payDay:       12,
			wantResponse: Response{Message: "Invalid how-much url"},
		},
		{
			name:         "Invalid pay_day parameter",
			methodType:   http.MethodGet,
			url:          "/till-salary/how-much?pay_day=32",
			statusInput:  http.StatusBadRequest,
			payDay:       32,
			wantResponse: Response{Message: "Invalid pay_day parameter"},
		},
		{
			name:         "Post method",
			methodType:   http.MethodPost,
			url:          "/till-salary/how-much?pay_day=15",
			statusInput:  http.StatusMethodNotAllowed,
			payDay:       15,
			wantResponse: Response{Message: "Method not allowed"},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			//create a new request with URL
			req, err := http.NewRequest(scenario.methodType, scenario.url, nil)
			if err != nil {
				t.Errorf("Error while creating new request: %v", err)
			}

			responseRecorder := httptest.NewRecorder()

			//call TillSalaryHandler with the new request and response recorder
			handler := http.HandlerFunc(TillSalaryHandler)
			handler.ServeHTTP(responseRecorder, req)

			//check status code
			gotStatus := responseRecorder.Code
			if gotStatus != scenario.statusInput {
				t.Errorf("Handler returned: %v, but i wanted: %v", gotStatus, http.StatusCreated)
			}

			var gotResponse Response
			err = json.Unmarshal(responseRecorder.Body.Bytes(), &gotResponse)
			if err != nil {
				t.Errorf("Error while unmarshaling response body: %v", err)
			}

			diff := cmp.Diff(gotResponse, scenario.wantResponse)
			if diff != "" {
				t.Errorf("Wanted body: %#v, got body: %#v\n%v", scenario.wantResponse, gotResponse, diff)
			}

		})
	}
}

func TestListPayDayDatesHandler(t *testing.T) {
	type testCases struct {
		name         string
		methodType   string
		url          string
		statusInput  int
		wantResponse Response
	}

	for _, scenario := range []testCases{
		{
			name:        "valid URL (15)",
			methodType:  http.MethodGet,
			url:         "/till-salary/pay-day/15/list-dates",
			statusInput: http.StatusCreated,
			wantResponse: Response{Data: map[string]interface{}{
				"pay_day_dates": []any{time.Date(2023, 3, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"), time.Date(2023, 4, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
					time.Date(2023, 5, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"), time.Date(2023, 6, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
					time.Date(2023, 7, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"), time.Date(2023, 8, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
					time.Date(2023, 9, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"), time.Date(2023, 10, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
					time.Date(2023, 11, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02"), time.Date(2023, 12, 15, 0, 0, 0, 0, time.Local).Format("2006-01-02")},
			}, Message: "Pay day dates"},
		},
		{
			name:         "invalid pay day parameter",
			methodType:   http.MethodGet,
			url:          "/till-salary/pay-day/32/list-dates",
			statusInput:  http.StatusBadRequest,
			wantResponse: Response{Message: "Invalid pay_day parameter"},
		},
		{
			name:         "invalid URL: list-dates is missing",
			methodType:   http.MethodGet,
			url:          "/till-salary/pay-day/14",
			statusInput:  http.StatusBadRequest,
			wantResponse: Response{Message: "Invalid pay-day url"},
		},
		{
			name:         "Post method",
			methodType:   http.MethodPost,
			url:          "/till-salary/pay-day/15/list-dates",
			statusInput:  http.StatusMethodNotAllowed,
			wantResponse: Response{Message: "Method not allowed"},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			// Create new request with url
			req, err := http.NewRequest(scenario.methodType, scenario.url, nil)
			if err != nil {
				t.Errorf("Error while creating new request: %#v", req)
			}

			responseRecorder := httptest.NewRecorder()

			// Call the ListPayDayDatesHandler function with the request and ResponseRecorder
			handler := http.HandlerFunc(ListPayDayDatesHandler)
			handler.ServeHTTP(responseRecorder, req)

			// Check the status code of the response
			status := responseRecorder.Code
			if status != scenario.statusInput {
				t.Errorf("handler returned wrong status code: got %#v want %#v", status, http.StatusCreated)
			}

			var gotResponse Response
			err = json.Unmarshal(responseRecorder.Body.Bytes(), &gotResponse)
			if err != nil {
				t.Errorf("Error while unmarshaling response body: %v", err)
			}

			diff := cmp.Diff(gotResponse, scenario.wantResponse)
			if diff != "" {
				t.Errorf("Wanted body: %#v, got body: %#v\n%v", scenario.wantResponse, gotResponse, diff)
			}
		})
	}
}
