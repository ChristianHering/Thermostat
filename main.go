package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/ChristianHering/Thermostat/utils"
	waveshare "github.com/ChristianHering/WaveShare"
	"golang.org/x/exp/io/i2c"
)

var i2cDevice *i2c.Device

func main() {
	err := utils.SetupUtils()
	if err != nil {
		panic(err)
	}

	i2cDevice, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x48)
	if err != nil {
		panic(err)
	}

	var catchFunc func()
	catchFunc = func() {
		defer func() {
			if err := recover(); err != nil {
				utils.LogError(fmt.Sprintf("%+v", err))

				catchFunc()
			}
		}()

		displayLoop() //Recovers from panics, then restarts displayLoop
	}
	catchFunc()
}

func displayLoop() {
	var sensorTemps []float64 = append([]float64{}, getTemp())
	var rData Data

	for true {
		weatherData, err := getWeatherData()
		if err != nil {
			panic(err)
		}

		//This is the number of times our program refreshes
		//local data before calling OpenWeather Map's API again.
		//
		//According to their website, OpenWeather Map updates data
		//for a given location no more than once every 10 minutes.
		for i := 0; i < 5; i++ {
			t, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp") //CPU temp readout
			if err != nil {
				panic(err)
			}

			temp, err := strconv.Atoi(strings.TrimSuffix(string(t), "\n"))
			if err != nil {
				panic(err)
			}

			t, err = ioutil.ReadFile("/proc/loadavg") //CPU load readout
			if err != nil {
				panic(err)
			}

			rData.RenderTime = time.Now().Format("3:04pm")
			rData.CPUTemp = strconv.FormatFloat(float64(temp/1000), 'f', -1, 64)
			rData.CPUUsage = strings.Split(string(t), " ")
			rData.SensorTemp = mean(sensorTemps)

			img := render(weatherData, rData)

			waveshare.Initialize()

			waveshare.DisplayImage(img)

			waveshare.Sleep()

			for i := 0; i < 10; i++ { //Average 10 temperature readings per minute to smooth the result out
				sensorTemps = append(sensorTemps, getTemp())

				if len(sensorTemps) > 10 {
					sensorTemps = sensorTemps[1:]
				}

				time.Sleep(6 * time.Second)
			}
		}
	}
}

//Returns the temperature from our TMP117 in fahrenheit
//
//Generally speaking, this function shouldn't be
//called more than once per second. For faster updates,
//more sensors should be used, sensor settings should
//be changed, or functions should cache readings.
func getTemp() float64 {
	var b = []byte{0, 0}

	err := i2cDevice.ReadReg(0x00, b)
	if err != nil {
		panic(err)
	}

	return (float64(binary.BigEndian.Uint16(b))*0.0078125)*1.8 + 32
}

//Calculates the average of a given array
func mean(arr []float64) float64 {
	var sum float64

	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}

	return sum / float64(len(arr))
}
