package main

import (
	"log"
	"math"
	"math/cmplx"
	"os"

	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	var reader *wav.Reader

	file, err := os.Open("./input/test.wav")
	if err != nil {
		log.Fatal(err)
	} else {
		reader = wav.NewReader(file)
	}
	defer file.Close()

	// srate, _ := reader.Samprate()
	// bit, _ := reader.BitsPerSample()
	d, _ := reader.Datasize()
	samples, _ := reader.ReadSamples(d)
	dataL := make([]float64, 0, len(samples))
	for _, sample := range samples {
		dataL = append(dataL, float64(reader.IntValue(sample, 0)))
	} //sampleをfloatに変換してdataLに追加

	specLcmplx := fft.FFTReal(dataL)
	specLcmplx = specLcmplx[:len(specLcmplx)/2]
	specL := []float64{}
	var max float64

	for _, bin := range specLcmplx {
		max = math.Max(max, plvtodb(cmplx.Abs(bin)))
		specL = append(specL, plvtodb(cmplx.Abs(bin)))
	} //fftした後のパワーdB値を格納

	for i, bin := range specL {
		specL[i] = bin - max
	}
	p := plot.New()
	// p.Title.Text = "Wavplot"
	// p.X.Label.Text = "Time"
	// p.Y.Label.Text = "Amplitude"

	p.Title.Text = "Spectrum"
	p.X.Label.Text = "Frequency(Hz)"
	p.Y.Label.Text = "Power(dB)"
	p.X.Scale = plot.LogScale{}

	// plotdata := datatoplt(dataL)
	// err1 := plotutil.AddLinePoints(p, "L", plotdata)
	// if err1 != nil {
	// 	panic(err1)
	// }
	// if err1 := p.Save(10*vg.Inch, 10*vg.Inch, "./output/wavplot.png"); err1 != nil {
	// 	panic(err1)
	// }

	l, err1 := plotter.NewLine(datatoplt(specL))
	if err1 != nil {
		panic(err1)
	}
	l.LineStyle.Width = vg.Points(1)

	p.Add(l)
	if err1 := p.Save(10*vg.Inch, 10*vg.Inch, "./output/line.png"); err != nil {
		panic(err1)
	}
}

func datatoplt(data []float64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i) + 1
		pts[i].Y = data[i]
	}
	return pts
}

func plvtodb(data float64) float64 {
	return 10 * math.Log10(math.Pow(data, 2))
}
