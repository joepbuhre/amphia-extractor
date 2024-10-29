package main

import (
	"encoding/json"
	"joepbuhre/amphia-agenda-ical/v2/models"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		isAllowedMethod := (func() bool {
			switch r.Method {
			case
				"GET",
				"POST":
				return true
			}
			return false
		})()

		if !isAllowedMethod {
			log.Printf("no allowed method, got [%s]", r.Method)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("origin"))
		w.Header().Set("Access-Control-Allow-Methods", r.Method)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		viper.SetConfigName(".env") // name of config file (without extension)
		viper.SetConfigType("env")  // REQUIRED if the config file does not have an extension
		viper.AutomaticEnv()
		viper.AddConfigPath(".") // optionally specify the path to look for the config file

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("error reading .env file: %s, this can make sense if you set it manually", err)
		}

		var config = models.Config{
			AgendaId:  viper.GetInt("AGENDA_ID"),
			AmphiaUrl: viper.GetString("AMPHIA_URL"),
			BaseUrl:   viper.GetString("BASE_URL"),
		}
		// Retrieve the values
		config.BearerToken = r.URL.Query().Get("bearer")
		shifts, err := models.MakeRequestWithBearerToken(config.AmphiaUrl, config.BearerToken)

		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)

		minDate, maxDate := getMinMaxDate(shifts)

		models.DeleteMeetingsInRange(config, minDate, maxDate)
		log.Printf("Deleted all meetings in range of currently fetched (%v to %v)", minDate, maxDate)

		for _, shift := range shifts {

			models.PostShiftToMeetings(config, "", shift)
			log.Printf("Posted shift to agenda with id %v", shift.Id)

			shiftStr, _ := json.Marshal(shift)

			w.Write(shiftStr)
			w.Write([]byte("\n"))

		}

	})

	log.Println("Server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getMinMaxDate(shifts []models.Shift) (minDate time.Time, maxDate time.Time) {
	var err error
	var beginDate time.Time = time.Now()
	var endDate time.Time = time.Now()

	minDate = beginDate
	maxDate = endDate

	for _, shift := range shifts {
		// Parse Begindate
		beginDate, err = time.Parse(time.RFC3339, shift.BeginDate)
		log.Println(beginDate)
		if err != nil {
			log.Fatal(err)
		}

		// Parse Enddate
		endDate, err = time.Parse(time.RFC3339, shift.EndDate)
		if err != nil {
			log.Fatal(err)
		}

		if beginDate.Before(minDate) {
			minDate = beginDate
		}
		if endDate.After(maxDate) {
			maxDate = endDate
		}
	}
	return
}
