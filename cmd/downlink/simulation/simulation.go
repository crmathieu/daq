package main
 
import (
    "fmt"
    "math"
)
 
type ypFunc func(t, y float64) float64
type ypStepFunc func(t, y, dt float64) float64
 
// newRKStep takes a function representing a differential equation
// and returns a function that performs a single step of the forth-order
// Runge-Kutta method.
func newRK4Step(yp ypFunc) ypStepFunc {
    return func(t, y, dt float64) float64 {
        dy1 := dt * yp(t, y)
        dy2 := dt * yp(t+dt/2, y+dy1/2)
        dy3 := dt * yp(t+dt/2, y+dy2/2)
        dy4 := dt * yp(t+dt, y+dy3)
        return y + (dy1+2*(dy2+dy3)+dy4)/6
    }
}
 
// example differential equation
func yprime(t, y float64) float64 {
    return t * math.Sqrt(y)
}
 
// exact solution of example
func actual(t float64) float64 {
    t = t*t + 4
    return t * t / 16
}
 
func main2() {
    t0, tFinal := 0, 10 // task specifies times as integers,
    dtPrint := 1        // and to print at whole numbers.
    y0 := 1.            // initial y.
    dtStep := .1        // step value.
 
    t, y := float64(t0), y0
    ypStep := newRK4Step(yprime)
    for t1 := t0 + dtPrint; t1 <= tFinal; t1 += dtPrint {
        printErr(t, y) // print intermediate result
        for steps := int(float64(dtPrint)/dtStep + .5); steps > 1; steps-- {
            y = ypStep(t, y, dtStep)
            t += dtStep
        }
        y = ypStep(t, y, float64(t1)-t) // adjust step to integer time
        t = float64(t1)
    }
    printErr(t, y) // print final result
}
 
func printErr(t, y float64) {
    fmt.Printf("y(%.1f) = %f Error: %e\n", t, y, math.Abs(actual(t)-y))
}
/*
Output:
y(0.0) = 1.000000 Error: 0.000000e+00
y(1.0) = 1.562500 Error: 1.457219e-07
y(2.0) = 3.999999 Error: 9.194792e-07
y(3.0) = 10.562497 Error: 2.909562e-06
y(4.0) = 24.999994 Error: 6.234909e-06
y(5.0) = 52.562489 Error: 1.081970e-05
y(6.0) = 99.999983 Error: 1.659460e-05
y(7.0) = 175.562476 Error: 2.351773e-05
y(8.0) = 288.999968 Error: 3.156520e-05
y(9.0) = 451.562459 Error: 4.072316e-05
y(10.0) = 675.999949 Error: 5.098329e-05
*/