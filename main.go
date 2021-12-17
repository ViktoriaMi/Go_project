package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	//"github.com/gin-gonic/gin"
	//"github.com/TheBookPeople/iso3166"
	"time"
)

// настройки для редиса
var ctx = context.Background()
var name_key = "weatherJsonStr"
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func RedisClientSet(body []byte) {
	//desc := dataW.Description
	//name_key := "weatherJsonStr"

	err := rdb.Set(ctx, name_key, body, 30*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func RedisClientGet() ([]byte, error) {
	val, err := rdb.Get(ctx, name_key).Bytes()
	if err != nil {
		//panic(err)
		var oshibka []byte
		return oshibka, err
	}
	//fmt.Println("key = ", val)
	return val, nil

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
}

type Weather struct {
	// Coord struct {
	// 	Lon float64 `json:"lon"`
	// 	Lat float64 `json:"lat"`
	// } `json:"coord"`
	Weather []struct {
		//ID          int    `json:"id"`
		//Main        string `json:"main"`
		Description string `json:"description"`
		//Icon        string `json:"icon"`
	} `json:"weather"`
	//Base string `json:"base"`
	Main struct {
		// температура в кельвинах
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		// давление в гектопаскалях
		Pressure float64 `json:"pressure"`
		// влажность, %
		Humidity int `json:"humidity"`
	} `json:"main"`
	//Visibility int `json:"visibility"`
	Wind struct {
		Speed float64 `json:"speed"`
		//Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		// облачность
		All int `json:"all"`
	} `json:"clouds"`
	//Dt  int `json:"dt"`
	Sys struct {
		//Type    int     `json:"type"`
		//ID      int     `json:"id"`
		//Message float64 `json:"message"`
		// код страны
		Country string `json:"country"`
		//Sunrise int     `json:"sunrise"`
		//Sunset  int     `json:"sunset"`
	} `json:"sys"`
	//Timezone int    `json:"timezone"`
	//ID       int    `json:"id"`
	// название города
	Name string `json:"name"`
	//Cod      int    `json:"cod"`
}

type Geo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

// получение ip
func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	//ips := r.Header.Get("X-FORWARDED-FOR")
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	x := fmt.Errorf("no valid ip found")
	return "", x
}

func weather(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nWeather processing")

	// ip, err := getIP(r)
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	w.Write([]byte("No valid ip"))
	// }
	// w.WriteHeader(200)
	// w.Write([]byte(ip))

	////////////////////////
	location := r.URL.Query().Get("location")
	if len(location) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"Location parameter is requiered\"}"))
		return
	}

	//str := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=%s",
	//location, api_key, "ru")

	str := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=%s&units=metric",
		location, os.Getenv("APIKEY"), "ru")

	resp, err := http.Get(str)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// используем редис !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//RedisClientSet(body)
	new_body, new_err := RedisClientGet()
	if new_err != nil {

		var data Weather
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
		}

		// перевод из гектопаскалей в мм. рт. столба
		data.Main.Pressure = data.Main.Pressure * 0.750064

		weathJSON, err := json.MarshalIndent(data, "", "   ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("First request / ttl expired\n")
		fmt.Printf("Weather data by city name:\n%s\n", string(weathJSON))

		RedisClientSet(body)

		//log.Fatalln(new_err)
	} else {
		var data Weather
		err = json.Unmarshal(new_body, &data)
		if err != nil {
			log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
		}

		// перевод из гектопаскалей в мм. рт. столба
		data.Main.Pressure = data.Main.Pressure * 0.750064

		weathJSON, err := json.MarshalIndent(data, "", "   ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Weather data by city name from cache:\n%s\n", string(weathJSON))
	}

	// ответ строкой
	//response := string(body)
	// печать в консоль
	//log.Println(response)
	//log.Println(time.Now().UTC())
	//bodybyte := [byte]new_body

	//// 2 ЛАБА
	// отправляем запрос к серверу, чтобы он определил ip
	str2 := "http://ip-api.com/json/"

	resp2, err2 := http.Get(str2)
	if err2 != nil {
		log.Fatalln(err2)
	}
	body2, err2 := ioutil.ReadAll(resp2.Body)
	if err2 != nil {
		log.Fatalln(err2)
	}

	var data2 Geo
	err2 = json.Unmarshal(body2, &data2)
	if err2 != nil {
		log.Fatalf("Error occured during unmarshaling 2. Error: %s", err2.Error())
	}

	//geoJSON, err2 := json.MarshalIndent(data2, "", "   ")
	//if err2 != nil {
	//log.Fatalf(err2.Error())
	//}
	//fmt.Printf("Weather data2:\n%s\n", string(geoJSON))

	ip_str := data2.Query
	city_str := data2.RegionName

	fmt.Printf("\nIP = %s\n", ip_str)

	//отправляем запрос openweather для определения погоды по ip

	str3 := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=%s&units=metric",
		city_str, os.Getenv("APIKEY"), "ru")

	resp3, err3 := http.Get(str3)
	if err3 != nil {
		log.Fatalln(err3)
	}
	body3, err3 := ioutil.ReadAll(resp3.Body)
	if err3 != nil {
		log.Fatalln(err3)
	}

	var data3 Weather
	err3 = json.Unmarshal(body3, &data3)
	if err3 != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err3.Error())
	}

	data3.Main.Pressure = data3.Main.Pressure * 0.750064

	weathJSON2, err3 := json.MarshalIndent(data3, "", "   ")
	if err3 != nil {
		log.Fatalf(err3.Error())
	}
	fmt.Printf("Weather data in %s by IP:\n%s\n", city_str, string(weathJSON2))
}

func main() {
	server := http.Server{
		//Addr: "localhost:6379",
		Addr: "0.0.0.0:8080",
	}

	http.HandleFunc("/weather/", weather)
	server.ListenAndServe()
}
