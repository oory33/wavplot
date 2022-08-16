package main

import (
	"log"
	"os"

	"github.com/oory33/wavplot/plt"

	"github.com/youpy/go-wav"
)

func main() {
	var reader *wav.Reader

	file, err := os.Open("./input/H4.wav")
	if err != nil {
		log.Fatal(err)
	} else {
		reader = wav.NewReader(file)
	}
	defer file.Close()

	plt.Spec(reader, 0, "freq")
}
