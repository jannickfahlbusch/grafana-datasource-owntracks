package datasource

import (
	"testing"

	"github.com/jannickfahlbusch/owntracks-go/types"
)

func generateLocations(i int) *types.LocationList {
	locationList := &types.LocationList{}
	locationList.Count = i

	for n := 0; n < i; n++ {
		locationList.Data = append(locationList.Data, types.Location{
			Longitude: float64(i),
			Latitude:  float64(i),
		})
	}

	return locationList
}

func benchmark_toTimeSeries(i int, b *testing.B) {
	locationList := generateLocations(i)

	for n := 0; n < b.N; n++ {
		_, err := toTimeSeries(locationList)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_toTimeSeries100(b *testing.B)      { benchmark_toTimeSeries(100, b) }
func Benchmark_toTimeSeries1000(b *testing.B)     { benchmark_toTimeSeries(1000, b) }
func Benchmark_toTimeSeries10000(b *testing.B)    { benchmark_toTimeSeries(10000, b) }
func Benchmark_toTimeSeries100000(b *testing.B)   { benchmark_toTimeSeries(100000, b) }
func Benchmark_toTimeSeries1000000(b *testing.B)  { benchmark_toTimeSeries(1000000, b) }
func Benchmark_toTimeSeries10000000(b *testing.B) { benchmark_toTimeSeries(10000000, b) }

func benchmark_toTable(i int, b *testing.B) {
	locationList := generateLocations(i)

	for n := 0; n < b.N; n++ {
		_, err := toTable(locationList)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_toTable100(b *testing.B)      { benchmark_toTable(100, b) }
func Benchmark_toTable1000(b *testing.B)     { benchmark_toTable(1000, b) }
func Benchmark_toTable10000(b *testing.B)    { benchmark_toTable(10000, b) }
func Benchmark_toTable100000(b *testing.B)   { benchmark_toTable(100000, b) }
func Benchmark_toTable1000000(b *testing.B)  { benchmark_toTable(1000000, b) }
func Benchmark_toTable10000000(b *testing.B) { benchmark_toTable(10000000, b) }
