package interp

import "math"

const NumSamples = 128

var log1 = math.Log10(1.41)
var log3 = math.Log10(3)

func SampleLog1(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * math.Log10(1 + 0.41 * x) / log1))
	}
	return pts
}

func SampleLog3(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * math.Log10(1 + 2 * x) / log3))
	}
	return pts
}

func SampleSine(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * ((1.0 - math.Cos(math.Pi * ((x - x1) / (x2 - x1)))) / 2.0)))
	}
	return pts
}

func SampleSCurve(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * (3.0 * math.Sqrt(x) - 2 * math.Pow(x, 3.0))))
	}
	return pts
}

func SampleInvertSCurve(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * ((1 - (3.0 * math.Sqrt(x) - 2 * math.Pow(x, 3.0))))))
	}
	return pts
}

func SampleLinear(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * ((x - x1) / (x2 - x1))))
	}
	return pts
}

func SampleExp1(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * (math.Pow(1.41, x) - 1) / 0.41))
	}
	return pts
}

func SampleExp3(x1 float64, y1 float64, x2 float64, y2 float64) []float32 {
	pts := make([]float32, 0, NumSamples)
	deltaY := y2 - y1
	for x := x1; x < x2; x += (x2 - x1) / NumSamples {
		pts = append(pts, float32(y1 + deltaY * (math.Pow(3, x) - 1) / 2))
	}
	return pts
}

func SampleConst(y float32) []float32 {
	pts := make([]float32, 0, NumSamples)
	for range NumSamples {
		pts = append(pts, y)
	}
	return pts
}
