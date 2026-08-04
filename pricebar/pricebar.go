package pricebar

import (
	"fmt"
	"math"
	"time"
)

type Config struct {
	Symbol     string
	Date       uint64
	Time       uint32
	BasisPrice float64
	Timestamp  time.Time
}

type PriceBar struct {
	// binary properties saved to disk
	Id     uint64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Vwap   float64
	Volume uint64

	// additional properties not saved to disk
	Symbol string
	Date   uint64
	Time   uint32
}

func GenerateId(date uint64, time uint32) uint64 {
	return (date << 16) + uint64(time)
}

func ExtractDateFromPriceBarId(priceBarId uint64) uint64 {
	return priceBarId >> 16
}

func ExtractTimeFromPriceBarId(priceBarId uint64) uint32 {
	return uint32(priceBarId % (1 << 16))
}

func New(config Config) *PriceBar {
	bar := PriceBar{}
	bar.Symbol = config.Symbol
	bar.Date = config.Date
	bar.Time = config.Time

	// if timestamp provided, use that to derive date/time
	if !config.Timestamp.IsZero() {
		dt := (config.Timestamp.Year() * 10000) + int(config.Timestamp.Month()*100) + config.Timestamp.Day()
		tm := (config.Timestamp.Hour() * 100) + config.Timestamp.Minute()
		bar.Date = uint64(dt)
		bar.Time = uint32(tm)
	}

	if bar.Date != 0 {
		bar.Id = GenerateId(bar.Date, bar.Time)
	}

	if config.BasisPrice != 0.0 {
		bar.Open = config.BasisPrice
		bar.High = config.BasisPrice
		bar.Low = config.BasisPrice
		bar.Close = config.BasisPrice
		bar.Vwap = config.BasisPrice
	}

	return &bar
}

func (bar *PriceBar) IsUp() bool {
	return bar.IsUpBy(0.0)
}

func (bar *PriceBar) IsUpBy(threshold float32) bool {
	return bar.Close > bar.Open+float64(threshold)
}

func (bar *PriceBar) IsDown() bool {
	return bar.IsDownBy(0.0)
}

func (bar *PriceBar) IsDownBy(threshold float32) bool {
	return bar.Close < bar.Open-float64(threshold)
}

func (bar *PriceBar) GetHour() uint32 {
	return bar.Time / 100
}

func (bar *PriceBar) GetMinute() uint32 {
	return bar.Time % 100
}

func (bar *PriceBar) GetRange() float64 {
	return bar.High - bar.Low
}

func (bar *PriceBar) GetBody() float64 {
	return math.Abs(bar.Close - bar.Open)
}

func (bar *PriceBar) GetHighWick() float64 {
	return bar.High - math.Max(bar.Open, bar.Close)
}

func (bar *PriceBar) GetLowWick() float64 {
	return math.Min(bar.Open, bar.Close) - bar.Low
}

func (bar *PriceBar) FillPriceDataFrom(other *PriceBar) {
	bar.Open = other.Open
	bar.High = other.High
	bar.Low = other.Low
	bar.Close = other.Close
	bar.Vwap = other.Vwap
	bar.Volume = other.Volume
}

func (bar *PriceBar) Aggregate(high float64, low float64, close float64, volume uint64, vwap float64) {
	bar.High = math.Max(bar.High, high)
	bar.Low = math.Min(bar.Low, low)
	bar.Close = close

	totalVolume := bar.Volume + volume
	currentTotal := bar.Vwap * float64(bar.Volume)
	updateTotal := vwap * float64(volume)

	bar.Volume = totalVolume
	if totalVolume > 0 {
		bar.Vwap = (currentTotal + updateTotal) / float64(totalVolume)
	}
}

func (bar *PriceBar) ToString(decimalPlaces int) string {
	decimalFormat := fmt.Sprintf("%%.%df", decimalPlaces)
	outputFormat := fmt.Sprintf("%%d %%4d O=%[1]s H=%[1]s L=%[1]s C=%[1]s V=%%d", decimalFormat)
	return fmt.Sprintf(outputFormat, bar.Date, bar.Time, bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
}
