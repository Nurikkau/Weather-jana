package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

// JSON struct to hold the weather data
type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

// get weather data from the API
func getWeather(city string, apiKey string) (WeatherData, error) {
	// Units=metric преобразует Фаренгейты в Цельсии
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("ошибка API: %d", resp.StatusCode)
	}

	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherData{}, err
	}
	return data, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		city := r.FormValue("city")
		apiKey := "ТВОЙ_КЛЮЧ_ТУТ" // В идеале брать из os.Getenv("API_KEY")

		data, err := getWeather(city, apiKey)
		if err != nil {
			// Передаем ошибку в шаблон, чтобы пользователь её увидел
			tmpl.Execute(w, map[string]string{"Error": "Город не найден!"})
			return
		}
		tmpl.Execute(w, data)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", weatherHandler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
