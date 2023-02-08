package main

import (
	"fmt"
	"os"
)

func main() {

	imageHeight := 256
	imageWidth := 256

	write(fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight))

	for j := imageHeight - 1; j >= 0; j-- {
		log(fmt.Sprintf("\rScanlines remaining: %d", j))
		for i := 0; i < imageWidth; i++ {
			r := float64(i) / float64(imageWidth-1)
			g := float64(j) / float64(imageHeight-1)
			b := 0.25

			ir := int(255.999 * r)
			ig := int(255.999 * g)
			ib := int(255.999 * b)

			write(fmt.Sprintf("%d %d %d\n", ir, ig, ib))
		}
	}

	log("\nDone")
}

func write(s string) {
	fmt.Printf(s)
}

func log(s string) {
	os.Stderr.WriteString(s)
}
