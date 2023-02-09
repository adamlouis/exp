package main

import (
	"fmt"
	"os"
)

type Color = Vec3
type Point3 = Vec3

func main() {

	// Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)

	// Camera
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	focalLength := 1.0

	origin := &Point3{0, 0, 0}
	horizontal := &Vec3{viewportWidth, 0, 0}
	vertical := &Vec3{0, viewportHeight, 0}

	lowerLeftCorner := origin.Sub(horizontal.Div(2)).Sub(vertical.Div(2)).Sub(&Vec3{0, 0, focalLength})

	write(fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight))

	for j := imageHeight - 1; j >= 0; j-- {
		log(fmt.Sprintf("\rScanlines remaining: %d", j))
		for i := 0; i < imageWidth; i++ {
			u := float64(i) / (float64(imageWidth) - 1)
			v := float64(j) / (float64(imageHeight) - 1)
			r := &Ray{origin, lowerLeftCorner.Add(horizontal.Mult(u)).Add(vertical.Mult(v)).Sub(origin)}
			pixelColor := getRayColor(r)
			writeColor(pixelColor)
		}
	}

	log("\nDone\n")
}

func getRayColor(r *Ray) *Color {

	if (hitSphere(&Point3{0, 0, -1}, 0.5, r)) {
		return &Color{1, 0, 0}
	}

	unitDirection := r.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	return (&Color{1.0, 1.0, 1.0}).Mult(1.0 - t).Add((&Color{0.5, 0.7, 1.0}).Mult(t))
}

func hitSphere(center *Point3, radius float64, r *Ray) bool {
	oc := r.Origin.Sub(center)

	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}

func write(s string) {
	fmt.Print(s)
}

func writeColor(c *Color) {
	write(fmt.Sprintf("%d %d %d\n", int(255.999*c.X), int(255.999*c.Y), int(255.999*c.Z)))
}

func log(s string) {
	os.Stderr.WriteString(s)
}
