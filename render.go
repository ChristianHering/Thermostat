package main

import (
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/fogleman/gg"

	"github.com/ChristianHering/Thermostat/utils"
)

const (
	width  = 800
	height = 480
)

//Data holds extra information rendered by our renderer function
type Data struct {
	RenderTime string
	CPUUsage   []string
	CPUTemp    string
	SensorTemp float64
}

//CardinalDirectionMap is for converting degrees to cardinal directions (for easy viewing)
var CardinalDirectionMap = map[int]string{
	0:  "N",
	1:  "NNE",
	2:  "NE",
	3:  "NEE",
	4:  "E",
	5:  "SEE",
	6:  "SE",
	7:  "SSE",
	8:  "S",
	9:  "SSW",
	10: "SW",
	11: "SWW",
	12: "W",
	13: "NWW",
	14: "NW",
	15: "NNW",
	16: "N",
}

//WeatherIDMap keeps track of all weather id associations
var WeatherIDMap = map[string]string{
	//----- DAY ASSOCIATIONS -----
	"200d": "wi-day-snow-thunderstorm", //2xx Group
	"201d": "wi-day-storm-showers",
	"202d": "wi-day-thunderstorm",
	"210d": "wi-day-lightning",
	"211d": "wi-day-lightning",
	"212d": "wi-day-lightning",
	"221d": "wi-day-lightning",
	"230d": "wi-day-snow-thunderstorm",
	"231d": "wi-day-storm-showers",
	"232d": "wi-day-thunderstorm",
	"300d": "wi-day-sprinkle", //3xx Group
	"301d": "wi-day-sprinkle",
	"302d": "wi-day-showers",
	"310d": "wi-day-snow",
	"311d": "wi-day-sprinkle",
	"312d": "wi-day-showers",
	"313d": "wi-day-rain-mix",
	"314d": "wi-day-rain-mix",
	"321d": "wi-day-showers",
	"500d": "wi-day-sprinkle", //5xx Group
	"501d": "wi-day-showers",
	"502d": "wi-day-rain",
	"503d": "wi-day-rain",
	"504d": "wi-day-rain",
	"511d": "wi-day-rain",
	"520d": "wi-day-sprinkle",
	"521d": "wi-day-showers",
	"522d": "wi-day-rain",
	"531d": "wi-day-rain",
	"600d": "wi-day-snow", //6xx Group
	"601d": "wi-day-snow-wind",
	"602d": "wi-day-snow-thunderstorm",
	"611d": "wi-day-sleet-storm",
	"612d": "wi-day-sleet",
	"613d": "wi-day-sleet-storm",
	"615d": "wi-day-rain-mix",
	"616d": "wi-day-rain-mix",
	"620d": "wi-day-snow",
	"621d": "wi-day-snow-wind",
	"622d": "wi-day-snow-thunderstorm",
	"701d": "wi-day-haze", //7xx Group
	"711d": "wi-day-haze",
	"721d": "wi-day-haze",
	"731d": "wi-day-windy",
	"741d": "wi-day-fog",
	"751d": "wi-day-haze",
	"761d": "wi-day-haze",
	"762d": "wi-day-haze",
	"771d": "wi-day-lightning",
	"781d": "wi-tornado",
	"800d": "wi-day-sunny",  //Clear
	"801d": "wi-day-cloudy", //80x Group
	"802d": "wi-day-cloudy",
	"803d": "wi-day-cloudy",
	"804d": "wi-day-sunny-overcast",
	//----- NIGHT ASSOCIATIONS -----
	"200n": "wi-night-alt-snow-thunderstorm", //2xx Group
	"201n": "wi-night-alt-storm-showers",
	"202n": "wi-night-alt-thunderstorm",
	"210n": "wi-night-alt-lightning",
	"211n": "wi-night-alt-lightning",
	"212n": "wi-night-alt-lightning",
	"221n": "wi-night-alt-lightning",
	"230n": "wi-night-alt-snow-thunderstorm",
	"231n": "wi-night-alt-storm-showers",
	"232n": "wi-night-alt-thunderstorm",
	"300n": "wi-night-alt-sprinkle", //3xx Group
	"301n": "wi-night-alt-sprinkle",
	"302n": "wi-night-alt-showers",
	"310n": "wi-night-alt-snow",
	"311n": "wi-night-alt-sprinkle",
	"312n": "wi-night-alt-showers",
	"313n": "wi-night-alt-rain-mix",
	"314n": "wi-night-alt-rain-mix",
	"321n": "wi-night-alt-showers",
	"500n": "wi-night-alt-sprinkle", //5xx Group
	"501n": "wi-night-alt-showers",
	"502n": "wi-night-alt-rain",
	"503n": "wi-night-alt-rain",
	"504n": "wi-night-alt-rain",
	"511n": "wi-night-alt-rain",
	"520n": "wi-night-alt-sprinkle",
	"521n": "wi-night-alt-showers",
	"522n": "wi-night-alt-rain",
	"531n": "wi-night-alt-rain",
	"600n": "wi-night-alt-snow", //6xx Group
	"601n": "wi-night-alt-snow-wind",
	"602n": "wi-night-alt-snow-thunderstorm",
	"611n": "wi-night-alt-sleet-storm",
	"612n": "wi-night-alt-sleet",
	"613n": "wi-night-alt-sleet-storm",
	"615n": "wi-night-alt-rain-mix",
	"616n": "wi-night-alt-rain-mix",
	"620n": "wi-night-alt-snow",
	"621n": "wi-night-alt-snow-wind",
	"622n": "wi-night-alt-snow-thunderstorm",
	"701n": "wi-night-fog", //7xx Group
	"711n": "wi-night-fog",
	"721n": "wi-night-fog",
	"731n": "wi-night-alt-cloudy-gusts",
	"741n": "wi-night-fog",
	"751n": "wi-night-fog",
	"761n": "wi-night-fog",
	"762n": "wi-night-fog",
	"771n": "wi-night-alt-lightning",
	"781n": "wi-tornado",
	"800n": "wi-night-clear",      //Clear
	"801n": "wi-night-alt-cloudy", //80x Group
	"802n": "wi-night-alt-cloudy",
	"803n": "wi-night-alt-cloudy",
	"804n": "wi-night-alt-partly-cloudy",
}

//Render renders and returns a new image to
//display with the given weather data
func render(weatherData WeatherData, rData Data) image.Image {
	var sunTime time.Time
	var sunState string

	for i := 0; sunTime.Before(time.Now()); i++ {
		if time.Unix(weatherData.Daily[i].SunriseTime, 0).After(time.Now()) {
			sunTime = time.Unix(weatherData.Daily[i].SunriseTime, 0)
			sunState = "Sunrise"
			break
		} else if time.Unix(weatherData.Daily[i].SunsetTime, 0).After(time.Now()) {
			sunTime = time.Unix(weatherData.Daily[i].SunsetTime, 0)
			sunState = "Sunset"
			break
		}
	}

	img, err := gg.LoadPNG("./render/background.png")
	if err != nil {
		panic(err)
	}

	ctx := gg.NewContextForImage(img)

	//Text rendering to background image

	ctx.SetColor(color.RGBA{0, 0, 0, 255})

	ctx.LoadFontFace("./render/fonts/OpenSans-Bold.ttf", 100)

	ctx.DrawStringAnchored(strconv.FormatFloat(float64(weatherData.Current.Temperature), 'f', 0, 32)+"°", 232.0, 100.0, 0.5, 0) //Outdoor Temperature
	ctx.DrawStringAnchored(strconv.FormatFloat(rData.SensorTemp, 'f', 1, 32)+"°", 362.0, 100.0, 0.0, 0)                         //Indoor Temperature TODO
	//ctx.DrawStringAnchored("72°", 784.0, 100.0, 1.0, 0)                                                                               //Used for setting indoor temperature

	ctx.LoadFontFace("./render/fonts/OpenSans-Bold.ttf", 40)

	//Left Data Panel
	ctx.DrawStringAnchored("Feels Like: "+strconv.FormatFloat(float64(weatherData.Current.TemperaturePerception), 'f', 0, 32)+"°", 175.0, 150.0, 0.5, 0) //Spelling of "Real" intentional
	ctx.DrawStringAnchored("Humidity: "+strconv.Itoa(weatherData.Current.Humidity)+"%", 175.0, 200.0, 0.5, 0)
	ctx.DrawStringAnchored("PoP (rain): "+strconv.Itoa(int(weatherData.Daily[0].Pop*100))+"%", 175.0, 250.0, 0.5, 0)

	//Right Data Panel
	ctx.DrawStringAnchored("Wind: "+CardinalDirectionMap[int(float64(weatherData.Current.WindDirection)/22.5)]+" @ "+strconv.Itoa(int(weatherData.Current.WindSpeed))+"mph", 570.5, 150.0, 0.5, 0)
	ctx.DrawStringAnchored(sunState+" at "+sunTime.Format("3:04pm"), 570.5, 200.0, 0.5, 0)
	ctx.DrawStringAnchored("Pressure: "+strconv.FormatFloat(float64(weatherData.Current.AtmosphericPressure)*0.02953, 'f', 2, 32)+"inHg", 570.5, 250.0, 0.5, 0)

	ctx.LoadFontFace("./render/fonts/OpenSans-Bold.ttf", 30)

	var p = [][]float64{} //Positional data for text rasterization

	p[0] = []float64{88.0, 305.0}
	p[1] = []float64{245.0, 305.0}
	p[2] = []float64{397.0, 305.0}
	p[3] = []float64{550.0, 305.0}
	p[4] = []float64{709.0, 305.0}

	for i := 0; i < 4; i++ { //5 day weekday labels
		ctx.DrawStringAnchored(time.Unix(weatherData.Daily[i+1].UnixTime, 0).Format("Mon"), p[i][0], p[i][1], 0.5, 0.5)
	}

	ctx.LoadFontFace("./render/fonts/OpenSans-Bold.ttf", 25)

	p = [][]float64{}

	p[0] = []float64{24.0, 450.0}
	p[1] = []float64{153.0, 450.0}

	p[2] = []float64{185.0, 450.0}
	p[3] = []float64{306.0, 450.0}

	p[4] = []float64{338.0, 450.0}
	p[5] = []float64{456.0, 450.0}

	p[6] = []float64{488.0, 450.0}
	p[7] = []float64{612.0, 450.0}

	p[8] = []float64{644.0, 450.0}
	p[9] = []float64{775.0, 450.0}

	for i := 0; i < 4; i++ { //5 day min/max forcast
		ctx.DrawStringAnchored(strconv.FormatFloat(float64(weatherData.Daily[i+1].Temperature.MinTemperature), 'f', 0, 32)+"°", p[i*2][0], p[i*2][1], 0.0, 0)
		ctx.DrawStringAnchored(strconv.FormatFloat(float64(weatherData.Daily[i+1].Temperature.MaxTemperature), 'f', 0, 32)+"°", p[(i*2)+1][0], p[(i*2)+1][1], 1.0, 0)
	}

	ctx.LoadFontFace("./render/fonts/OpenSans-Bold.ttf", 20)

	logStr := "Data Updated: " + rData.RenderTime + " -" //Build-out/render of log string
	logStr += " Adv CPU Load: " + rData.CPUUsage[0] + "% -"
	logStr += " CPU Temp: " + rData.CPUTemp + "°C -"
	logStr += " Errors: " + strconv.Itoa(utils.Errors)

	ctx.DrawStringAnchored(logStr, 400, 468, 0.5, 0.5)

	//Icon rendering to display image

	currentState := weatherData.Current.Weather[0].IconID
	currentState = currentState[len(currentState)-1:]

	img = getIconImage(WeatherIDMap[strconv.Itoa(weatherData.Current.Weather[0].ID)+currentState]) //Main Weather Icon
	ctx.DrawImageAnchored(img, 65, 65, 0.5, 0.5)

	p = [][]float64{}

	p[0] = []float64{69.0, 376.0}
	p[1] = []float64{245.0, 376.0}
	p[2] = []float64{397.0, 376.0}
	p[3] = []float64{550.0, 376.0}
	p[4] = []float64{709.0, 376.0}

	for i := 0; i < 4; i++ { //5 Day Forcast
		img = getIconImage(WeatherIDMap[strconv.Itoa(weatherData.Daily[i+1].Weather[0].ID)+currentState])
		ctx.DrawImageAnchored(img, int(p[i][0]), int(p[i][1]), 0.5, 0.5)
	}

	return ctx.Image()
}

func getIconImage(icon string) image.Image {
	img, err := gg.LoadPNG("./render/icons/" + icon + ".png")
	if err != nil {
		img, err = gg.LoadPNG("./render/icons/wi-na.png")
		if err != nil {
			panic(err)
		}
	}

	return img
}
