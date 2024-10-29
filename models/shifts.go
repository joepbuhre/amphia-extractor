package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Department struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Shift struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Remark      string     `json:"remark"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Department  Department `json:"department"`
	BeginDate   string     `json:"beginDate"`
	EndDate     string     `json:"endDate"`
}

type MeetingRequest struct {
	ID            int    `json:"id"` // Optional for update
	AgendaID      int    `json:"agenda_id"`
	Summary       string `json:"summary"`
	Description   string `json:"description"`
	StartDateTime string `json:"start_datetime"`
	EndDateTime   string `json:"end_datetime"`
	Location      string `json:"location"`
	Color         string `json:"color"`
}

func PostShiftToMeetings(config Config, token string, shift Shift) {
	log.Printf("Posting shift [%v]", shift)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var summaryStr string = shift.Description
	if summaryStr == "" {
		summaryStr = shift.Remark
	}

	if summaryStr == "" {
		summaryStr = shift.Name
	}

	var meetingRequest MeetingRequest = MeetingRequest{
		ID:            shift.Id,
		AgendaID:      config.AgendaId,
		Summary:       cases.Title(language.Dutch).String(summaryStr),
		Description:   shift.Name,
		StartDateTime: shift.BeginDate,
		EndDateTime:   shift.EndDate,
		Location:      shift.Department.Name,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(meetingRequest)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new HTTP request
	req, err := http.NewRequest("PUT", config.BaseUrl, &buf)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("received non-OK HTTP status: %s", resp.Status)

		body, err := io.ReadAll(io.Reader(resp.Body))
		if err != nil {
			log.Printf("error reading response body: %v", err)
		}
		log.Println(string(body[:]))
	}

}

// MakeRequestWithBearerToken sends a GET request to the given URL with the Bearer token
// and parses the response into a Shift struct
func MakeRequestWithBearerToken(url string, token string) ([]Shift, error) {
	method := "GET"

	client := &http.Client{}

	now := time.Now()
	fromDate := now.AddDate(0, -1, 0).Format("2006-01-02")
	toDate := now.AddDate(0, 5, 0).Format("2006-01-02")

	url = fmt.Sprintf(url+"?fromDate=%s&untilDate=%s", fromDate, toDate)

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("tenant", "amphiazh")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.Reader(res.Body))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if res.StatusCode > 399 {
		return nil, fmt.Errorf("something went wrong %s", string(body))
	}

	// Parse the JSON response into the Shift struct
	var shifts []Shift
	err = json.Unmarshal(body, &shifts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return shifts, nil
}

// DeleteMeetingsInRange
func DeleteMeetingsInRange(config Config, fromDate time.Time, toDate time.Time) error {
	method := "DELETE"
	url := config.BaseUrl

	client := &http.Client{}

	fromDateStr := fromDate.Format("2006-01-02")
	toDateStr := toDate.Format("2006-01-02")

	url = fmt.Sprintf(url+"&from_date=%s&to_date=%s", fromDateStr, toDateStr)

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.Reader(res.Body))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.StatusCode > 399 {
		return fmt.Errorf("something went wrong %s", string(body))
	}

	return nil
}
