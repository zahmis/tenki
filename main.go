package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"tenki.com/tenki/ps"
)

func main() {
	weatherApiToken := os.Getenv("WEATHER_API_TOKEN")

	if weatherApiToken == "" {
		log.Fatal("WEATHER_API_TOKEN is not set")
		return
	}

    // OpenWeatherMap APIから天気情報を取得するためのURL
    url := "http://api.openweathermap.org/data/2.5/weather?q=Tokyo&units=metric&appid=" + weatherApiToken

    // HTTP GETリクエストを作成
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    // HTTP GETリクエストを送信
    client := new(http.Client)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    // HTTPレスポンスを読み込む
    defer resp.Body.Close()

 

    // JSONをパース
    var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		log.Fatal(err)
	}

	weather := ps.Weather{
		Main: weatherData["weather"].([]interface{})[0].(map[string]interface{})["main"].(string),
		Description: weatherData["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
	}

	main := ps.Main{
		Temp: weatherData["main"].(map[string]interface{})["temp"].(float64),
		Pressure: weatherData["main"].(map[string]interface{})["pressure"].(float64),
		Humidity: weatherData["main"].(map[string]interface{})["humidity"].(float64),
		TempMin: weatherData["main"].(map[string]interface{})["temp_min"].(float64),
		TempMax: weatherData["main"].(map[string]interface{})["temp_max"].(float64),
	}

	wind := ps.Wind{
		Speed: weatherData["wind"].(map[string]interface{})["speed"].(float64),
		Deg: weatherData["wind"].(map[string]interface{})["deg"].(float64),
	}

	clouds := ps.Clouds{
		All: weatherData["clouds"].(map[string]interface{})["all"].(float64),
	}

	sys := ps.Sys{
		Country: weatherData["sys"].(map[string]interface{})["country"].(string),
		Sunrise: weatherData["sys"].(map[string]interface{})["sunrise"].(float64),
		Sunset: weatherData["sys"].(map[string]interface{})["sunset"].(float64),
	}

	totalInfo := ps.TotalInfo{
		Name: weatherData["name"].(string),
		Weather: weather,
		Main: main,
		Wind: wind,
		Clouds: clouds,
		Sys: sys,
	}

	fmt.Println(totalInfo)

	// JSONをパース
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	// 天気情報を表示
	fmt.Println("都市名:", data["name"])
	fmt.Println("天気:", data["weather"].([]interface{})[0].(map[string]interface{})["main"])
	fmt.Println("天気詳細:", data["weather"].([]interface{})[0].(map[string]interface{})["description"])
	fmt.Println("最低気温:", data["main"].(map[string]interface{})["temp_min"])
	fmt.Println("最高気温:", data["main"].(map[string]interface{})["temp_max"])
	fmt.Println("湿度:", data["main"].(map[string]interface{})["humidity"])
	fmt.Println("風速:", data["wind"].(map[string]interface{})["speed"])
	fmt.Println("雲量:", data["clouds"].(map[string]interface{})["all"])
	fmt.Println("国名:", data["sys"].(map[string]interface{})["country"])
	fmt.Println("日の出:", data["sys"].(map[string]interface{})["sunrise"])
	fmt.Println("日の入り:", data["sys"].(map[string]interface{})["sunset"])

	connectStr := "user=sizmayosimaz dbname=tenki2 sslmode=disable"

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブルを作成
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS totalInfo (
		id SERIAL PRIMARY KEY,
		name TEXT,
		weather TEXT,
		main TEXT,
		wind TEXT,
		clouds TEXT,
		sys TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS weather (
		id SERIAL PRIMARY KEY,
		main TEXT,
		description TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS main (
		id SERIAL PRIMARY KEY,
		temp TEXT,
		pressure TEXT,
		humidity TEXT,
		temp_min TEXT,
		temp_max TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS wind (
		id SERIAL PRIMARY KEY,
		speed TEXT,
		deg TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS clouds (
		id SERIAL PRIMARY KEY,
		all TEXT
	)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS sys (
		id SERIAL PRIMARY KEY,
		country TEXT,
		sunrise TEXT,
		sunset TEXT
	)`); err != nil {
		log.Fatal(err)
	}



	// データを挿入
	if _, err := db.Exec(`INSERT INTO weather (name, weather, main, wind, clouds, sys) VALUES ($1, $2, $3, $4, $5, $6)`,
		data["name"], data["weather"].([]interface{})[0].(map[string]interface{})["main"], data["main"].(map[string]interface{})["temp"], data["wind"].(map[string]interface{})["speed"], data["clouds"].(map[string]interface{})["all"], data["sys"].(map[string]interface{})["country"]); err != nil {
		log.Fatal(err)
	}

	if _,err := db.Exec(`INSERT INRO weather (main, description) VALUES ($1, $2)`,
		data["weather"].([]interface{})[0].(map[string]interface{})["main"], data["weather"].([]interface{})[0].(map[string]interface{})["description"]); err != nil {
		log.Fatal(err)
	}

	if _,err := db.Exec(`INSERT INRO main (temp, pressure, humidity, temp_min, temp_max) VALUES ($1, $2, $3, $4, $5)`,
		data["main"].(map[string]interface{})["temp"], data["main"].(map[string]interface{})["pressure"], data["main"].(map[string]interface{})["humidity"], data["main"].(map[string]interface{})["temp_min"], data["main"].(map[string]interface{})["temp_max"]); err != nil {
		log.Fatal(err)
	}

	if _,err := db.Exec(`INSERT INRO wind (speed, deg) VALUES ($1, $2)`,
		data["wind"].(map[string]interface{})["speed"], data["wind"].(map[string]interface{})["deg"]); err != nil {
		log.Fatal(err)
	}

	if _,err := db.Exec(`INSERT INRO clouds (all) VALUES ($1)`,
		data["clouds"].(map[string]interface{})["all"]); err != nil {
		log.Fatal(err)
	}

	if _,err := db.Exec(`INSERT INRO sys (country, sunrise, sunset) VALUES ($1, $2, $3)`,
		data["sys"].(map[string]interface{})["country"], data["sys"].(map[string]interface{})["sunrise"], data["sys"].(map[string]interface{})["sunset"]); err != nil {
		log.Fatal(err)
	}


	// データを取得
	rows, err := db.Query(`SELECT * FROM totalInfo`)
	if err != nil {
		log.Fatal(err)
	}
	


	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var weather string
		var main string
		var wind string
		var clouds string
		var sys string
		if err := rows.Scan(&id, &name, &weather, &main, &wind, &clouds, &sys); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, weather, main, wind, clouds, sys)
	}
}
