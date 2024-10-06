package main

import (
	"encoding/json"
	"joepbuhre/amphia-agenda-ical/v2/models"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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
