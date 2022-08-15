package plt

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
	"github.com/youpy/go-wav"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Wav(r *wav.Reader, ch uint, name string) {
	length, _ := r.Datasize()
	datas, _ := r.ReadSamples(length)

	plotdata := make([]float64, 0, length)

	for _, sample := range datas {
		plotdata = append(plotdata, float64(r.IntValue(sample, ch)))
	}
	p := plot.New()
	p.Title.Text = "Wavplot"
	p.X.Label.Text = "Time(sec)"
	p.Y.Label.Text = "Amplitude"

	fs, _ := r.Samprate()
	bit, _ := r.BitsPerSample()
	l, err1 := plotter.NewLine(normtime(plotdata, fs, bit))
	if err1 != nil {
		panic(err1)
	}
	l.LineStyle.Width = vg.Points(1)

	p.Add(l)
	if err := p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("./output/%s.png", name)); err != nil {
		panic(err)
	}
}

func Spec(r *wav.Reader, ch uint, name string) {
	length, _ := r.Datasize()
	datas, _ := r.ReadSamples(length)

	plotdata := make([]float64, 0, length)

	for _, sample := range datas {
		plotdata = append(plotdata, float64(r.IntValue(sample, ch)))
	}
	specLcmplx := fft.FFTReal(plotdata)
	specLcmplx = specLcmplx[:len(specLcmplx)/2]
	specL := []float64{}
	var max float64

	for _, bin := range specLcmplx {
		val := plvtodb(cmplx.Abs(bin))
		max = math.Max(max, val)
		specL = append(specL, val)
	}
	for i, bin := range specL {
		specL[i] = bin - max
	} //最大値とのdB化

	p := plot.New()
	p.Title.Text = "Specplot"
	p.X.Label.Text = "Frequency(Hz)"
	p.Y.Label.Text = "Power(dB)"
	p.X.Scale = plot.LogScale{}
	p.X.Min = 20

	fs, _ := r.Samprate()
	bit, _ := r.BitsPerSample()
	l, err1 := plotter.NewLine(normfreq(specL, fs, bit))
	if err1 != nil {
		panic(err1)
	}
	l.LineStyle.Width = vg.Points(1)

	p.Add(l)
	if err := p.Save(10*vg.Inch, 10*vg.Inch, fmt.Sprintf("./output/%s.png", name)); err != nil {
		panic(err)
	}
}

func normtime(data []float64, srate uint32, bit uint16) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i) / float64(srate)
		pts[i].Y = data[i] / math.Pow(2, float64(bit)-1)
	}
	return pts
}

func plvtodb(data float64) float64 {
	return 10 * math.Log10(math.Pow(data, 2))
}

func normfreq(data []float64, srate uint32, bit uint16) plotter.XYs {
	fn := float64(srate / 2)
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = ((float64(i) + 1) / float64(len(data))) * fn
		pts[i].Y = data[i]
	}
	return pts
}
