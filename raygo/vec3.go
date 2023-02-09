package main

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func NewZeroVec3() *Vec3 {
	return &Vec3{0, 0, 0}
}

func NewNegVec3(v *Vec3) *Vec3 {
	return &Vec3{-v.X, -v.Y, -v.Z}
}

func (v *Vec3) Add(v2 *Vec3) *Vec3 {
	// v.X += v2.X
	// v.Y += v2.Y
	// v.Z += v2.Z
	// return v
	return &Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v *Vec3) Sub(v2 *Vec3) *Vec3 {
	// v.X -= v2.X
	// v.Y -= v2.Y
	// v.Z -= v2.Z
	// return v
	return &Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v *Vec3) Mult(t float64) *Vec3 {
	// v.X *= t
	// v.Y *= t
	// v.Z *= t
	// return v
	return &Vec3{v.X * t, v.Y * t, v.Z * t}
}

func (v *Vec3) Div(t float64) *Vec3 {
	return v.Mult(1 / t)
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3) Printf() {
	fmt.Printf("[%.2f %.2f %.2f]\n", v.X, v.Y, v.Z)
}

func (v *Vec3) Unit() *Vec3 {
	return v.Div(v.Length())
}

func (v *Vec3) Dot(v2 *Vec3) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}
