package pricebar

import (
	"testing"
	"time"

	"github.com/achedges/go-assertions"
)

func TestPriceBar_GenerateId(t *testing.T) {
	var expectedId uint64 = 1325873235120
	var dt uint64 = 20231220
	var tm uint32 = 1200
	var testId = GenerateId(dt, tm)

	assertions.EqualUints(expectedId, testId, t)
}

func TestPriceBar_ExtractDateFromPriceBarId(t *testing.T) {
	var priceBarId uint64 = 1325873235120
	var expectedDate uint64 = 20231220
	var testDate = ExtractDateFromPriceBarId(priceBarId)

	assertions.EqualUints(expectedDate, testDate, t)
}

func TestPriceBar_ExtractTimeFromPriceBarId(t *testing.T) {
	var priceBarId uint64 = 1325873235120
	var expectedTime uint32 = 1200
	var testTime = ExtractTimeFromPriceBarId(priceBarId)

	assertions.EqualUints(expectedTime, testTime, t)
}

func TestPriceBar_New_WithSymbol(t *testing.T) {
	bar := New(Config{Symbol: "TEST"})

	expectedBar := PriceBar{
		Symbol: "TEST",
	}

	if *bar != expectedBar {
		t.Error()
	}
}

func TestPriceBar_New_WithSymbolDateTime(t *testing.T) {
	bar := New(Config{
		Symbol: "TEST",
		Date:   20260728,
		Time:   930,
	})

	expectedBar := PriceBar{
		Symbol: "TEST",
		Date:   20260728,
		Time:   930,
		Id:     GenerateId(20260728, 930),
	}

	if *bar != expectedBar {
		t.Error()
	}
}

func TestPriceBar_New_WithSymbolBasis(t *testing.T) {
	bar := New(Config{
		Symbol:     "TEST",
		BasisPrice: 10.0,
	})

	expectedBar := PriceBar{
		Symbol: "TEST",
		Open:   10.0,
		High:   10.0,
		Low:    10.0,
		Close:  10.0,
		Vwap:   10.0,
	}

	if *bar != expectedBar {
		t.Error()
	}
}

func TestPriceBar_New_WithSymbolDateTimeBasis(t *testing.T) {
	bar := New(Config{
		Symbol:     "TEST",
		Date:       20260728,
		Time:       930,
		BasisPrice: 10.0,
	})

	expectedBar := PriceBar{
		Symbol: "TEST",
		Date:   20260728,
		Time:   930,
		Id:     GenerateId(20260728, 930),
		Open:   10.0,
		High:   10.0,
		Low:    10.0,
		Close:  10.0,
		Vwap:   10.0,
	}

	if *bar != expectedBar {
		t.Error()
	}
}

func TestPriceBar_New_WithSymbolBasisTimestamp(t *testing.T) {
	bar := New(Config{
		Symbol:     "TEST",
		BasisPrice: 10.0,
		Timestamp:  time.Date(2026, 7, 28, 9, 30, 0, 0, time.UTC),
	})

	expectedBar := PriceBar{
		Symbol: "TEST",
		Date:   20260728,
		Time:   930,
		Id:     GenerateId(20260728, 930),
		Open:   10.0,
		High:   10.0,
		Low:    10.0,
		Close:  10.0,
		Vwap:   10.0,
	}

	if *bar != expectedBar {
		t.Error()
	}
}

func TestPriceBar_New_TimestampOverridesDateTime(t *testing.T) {
	bar := New(Config{
		Date:      20260701,
		Time:      930,
		Timestamp: time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC),
	})

	assertions.EqualUints(20260715, bar.Date, t)
	assertions.EqualUints(1200, bar.Time, t)
}

func TestPriceBar_New_NoDateNoId(t *testing.T) {
	bar := New(Config{
		Time: 930,
	})

	assertions.EqualUints(0, bar.Id, t)
}

func TestPriceBar_IsUp(t *testing.T) {
	upBar := PriceBar{
		Open:  10.0,
		Close: 11.0,
	}

	assertions.True(upBar.IsUp(), t)
	assertions.True(upBar.IsUpBy(0.9), t)
	assertions.False(upBar.IsUpBy(1.1), t)

	downBar := PriceBar{
		Open:  10.0,
		Close: 9.0,
	}

	assertions.False(downBar.IsUp(), t)
	assertions.False(downBar.IsUpBy(0.9), t)
	assertions.False(downBar.IsUpBy(1.1), t)
}

func TestPriceBar_IsDown(t *testing.T) {
	downBar := PriceBar{
		Open:  10.0,
		Close: 9.0,
	}

	assertions.True(downBar.IsDown(), t)
	assertions.True(downBar.IsDownBy(0.9), t)
	assertions.False(downBar.IsDownBy(1.1), t)

	upBar := PriceBar{
		Open:  10.0,
		Close: 11.0,
	}

	assertions.False(upBar.IsDown(), t)
	assertions.False(upBar.IsDownBy(0.9), t)
	assertions.False(upBar.IsDownBy(1.1), t)
}

func TestPriceBar_GetHour(t *testing.T) {
	bar := PriceBar{Time: 1440}
	assertions.EqualUints(14, bar.GetHour(), t)
}

func TestPriceBar_GetMinute(t *testing.T) {
	bar := PriceBar{Time: 1734}
	assertions.EqualUints(34, bar.GetMinute(), t)
}

func TestPriceBar_GetRange(t *testing.T) {
	bar := PriceBar{
		High: 10.0,
		Low:  9.0,
	}
	assertions.EqualFloats(1.0, bar.GetRange(), t)
}

func TestPriceBar_GetBody(t *testing.T) {
	upBar := PriceBar{
		Open:  10.0,
		Close: 11.0,
	}
	assertions.EqualFloats(1.0, upBar.GetBody(), t)

	downBar := PriceBar{
		Open:  10.0,
		Close: 9.0,
	}
	assertions.EqualFloats(1.0, downBar.GetBody(), t)
}

func TestPriceBar_GetHighWick(t *testing.T) {
	bar := PriceBar{
		Open:  10.0,
		Close: 11.0,
		High:  11.25,
	}
	assertions.EqualFloats(0.25, bar.GetHighWick(), t)
}

func TestPriceBar_GetLowWick(t *testing.T) {
	bar := PriceBar{
		Open:  10.0,
		Close: 11.0,
		Low:   9.75,
	}
	assertions.EqualFloats(0.25, bar.GetLowWick(), t)
}

func TestPriceBar_FillPriceDataFrom(t *testing.T) {
	bar := PriceBar{
		Open:   1.0,
		High:   2.0,
		Low:    0.9,
		Close:  1.5,
		Vwap:   1.35,
		Volume: 123,
	}

	var other = PriceBar{}
	other.FillPriceDataFrom(&bar)

	if other != bar {
		t.Error()
	}
}

func TestPriceBar_Aggregate(t *testing.T) {
	bar := PriceBar{
		Open:   10.0,
		High:   10.0,
		Low:    10.0,
		Close:  10.0,
		Volume: 100,
		Vwap:   10.0,
	}

	bar.Aggregate(10.5, 9.5, 9.75, 35, 10.2)
	assertions.EqualFloats(10.0, bar.Open, t)
	assertions.EqualFloats(10.5, bar.High, t)
	assertions.EqualFloats(9.5, bar.Low, t)
	assertions.EqualFloats(9.75, bar.Close, t)
	assertions.EqualUints(135, bar.Volume, t)
	assertions.CloseEnough(10.0518, bar.Vwap, 0.0001, t)
}

func TestPriceBar_ToString(t *testing.T) {
	bar := New(Config{
		Symbol:     "TEST",
		Date:       20260729,
		Time:       813,
		BasisPrice: 10.0,
	})
	assertions.EqualStrings("20260729  813 O=10.0 H=10.0 L=10.0 C=10.0 V=0", bar.ToString(1), t)

	bar.Time = 1234
	assertions.EqualStrings("20260729 1234 O=10.0 H=10.0 L=10.0 C=10.0 V=0", bar.ToString(1), t)
}
