package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type City struct {
	occurrences int
	minWeather float64
	averageWeather float64
	maxWeather float64
}

type CityMap map[string]City

func main() {
	file, err := os.Open("weather_stations.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	cities := make(CityMap)

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("Error reading csv file:", err)
		return
	}
	
	for _, item := range records {
		cityName, cityWeather := formatLine(item[0])

		if checkIfCityExistsInMap(cities, cityName) {
			updateCityInMap(&cities, cityName, cityWeather)
		} else {
			cities[cityName] = City{1, cityWeather, cityWeather, cityWeather}
		}
	}

	sortedCityNames := sortCitiesByName(cities)

	for _, cityName := range sortedCityNames {
		cityData := cities[cityName]
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", cityName, cityData.minWeather, cityData.averageWeather, cityData.maxWeather)
	}
}

func formatLine(text string) (cityName string, cityWeather float64) {
	stringSplit := strings.Split(text, ";")

	cityName = stringSplit[0]
	cityWeatherStr := stringSplit[1]
	cityWeather, err := strconv.ParseFloat(cityWeatherStr, 64)
	if err != nil {
		fmt.Println("Error converting weather data to float64: ", err)
		panic(err)
	}
	return cityName, cityWeather
}

func checkIfCityExistsInMap(cities CityMap, city string) bool {
	if _, exists := cities[city]; exists {
        return true
    }
    return false
}

func updateCityInMap(cities *CityMap, city string, weather float64) {
	currentCity := (*cities)[city]

	currentCity.occurrences += 1
	currentCity.averageWeather = (currentCity.averageWeather + weather) / float64(currentCity.occurrences)
	currentCity.minWeather = min(currentCity.minWeather, weather)
	currentCity.maxWeather = max(currentCity.maxWeather, weather)

	(*cities)[city] = currentCity
}

func sortCitiesByName(cities CityMap) []string {
	names := make([]string, 0, len(cities))
	for name := range cities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}