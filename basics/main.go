package main

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"golang.org/x/tour/pic"
)

type vertex struct {
	X, Y int
}

type fVertex struct {
	Lat, Long float64
}

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v[key]
}

var (
	v1 = vertex{1, 2}
	v2 = vertex{X: 1}
	v3 = vertex{}
	p  = &vertex{1, 2}
)

func sqrt(x float64) float64 {
	z := float64(1)

	for i := 1; i <= 5; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Printf("%d: %v\n", i, z)
	}

	return z
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	defer func() {
		fmt.Println("leaving main")
	}()

	{
		fmt.Println("--- switch")
		fmt.Print("Go runs on ")
		switch os := runtime.GOOS; os {
		case "darwin":
			fmt.Println("deniz.")
		case "linux":
			fmt.Println("Linux.")
		default:
			// freebsd, openbsd,
			// plan9, windows...
			fmt.Printf("%s.\n", os)
		}
	}

	{
		fmt.Println("--- loops")
		fmt.Println(sqrt(2), math.Sqrt(2))

		var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
		for i, v := range pow {
			fmt.Printf("2**%d = %d\n", i, v)
		}
	}

	{
		fmt.Println("--- structs")
		fmt.Println(v1, p, v2, v3)
	}

	{
		fmt.Println("--- arrays")
		var a [2]string
		a[0] = "Hello"
		a[1] = "array"
		fmt.Println(a[0], a[1])
		fmt.Println(a)

		primes := [6]int{2, 3, 5, 7, 11, 13}
		fmt.Println(primes)

		//slice
		var s []int = primes[1:4]
		fmt.Println(s)
	}

	{
		fmt.Println("--- slices are references to arrays")
		names := [4]string{
			"John",
			"Paul",
			"George",
			"Ringo",
		}
		fmt.Println(names)

		a := names[0:2]
		b := names[1:3]
		fmt.Println(a, b)

		b[0] = "XXX"
		fmt.Println(a, b)
		fmt.Println(names)

		// slice literals
		s := []int{2, 3, 5, 7, 11, 13}
		s = s[1:4]
		fmt.Println(s)

		s = s[:2]
		fmt.Println(s)

		s = s[1:]
		fmt.Println(s)
	}

	{
		fmt.Println("--- slice length and capacity")
		s := []int{2, 3, 5, 7, 11, 13}
		printSlice(s)

		// Slice the slice to give it zero length.
		s = s[:0]
		printSlice(s)

		// Extend its length.
		s = s[:4]
		printSlice(s)

		// Drop its first two values.
		s = s[2:]
		printSlice(s)

		var s2 []int
		printSlice(s2)
		if s2 == nil {
			fmt.Println("nil!")
		}

		s3 := make([]int, 0, 0)
		printSlice(s3)
	}

	{
		fmt.Println("--- slice append")
		var s []int
		printSlice(s)

		s = append(s, 0, 0, -3)
		printSlice(s)
	}

	{
		fmt.Println("--- slice of slice")
		pic.Show(func(dx, dy int) [][]uint8 {
			a := make([][]uint8, dy)
			for y := 0; y < dy; y++ {
				a[y] = make([]uint8, dx)
				for x := 0; x < dx; x++ {
					a[y][x] = uint8(x ^ y)
				}
			}
			return a
		})
	}

	{
		fmt.Println("--- maps")
		var m = map[string]fVertex{
			"Bell Labs": {40.68433, -74.39967},
			"Google":    {37.42202, -122.08408},
		}
		fmt.Println(m)

		n := make(map[string]int)
		n["Answer"] = 42
		fmt.Println("The value:", n["Answer"])

		n["Answer"] = 48
		fmt.Println("The value:", n["Answer"])

		delete(n, "Answer")
		fmt.Println("The value:", n["Answer"])

		v, ok := n["Answer"]
		fmt.Println("The value:", v, "Present?", ok)
	}

	{
		fmt.Println("--- interfaces")
		var i interface{} = "hello"
		s := i.(string)
		fmt.Println(s)

		s, ok := i.(string)
		fmt.Println(s, ok)

		f, ok := i.(float64)
		fmt.Println(f, ok)

		//f = i.(float64) // panic

		do(21)
		do("hello")
		do(true)
	}

	{
		fmt.Println("--- goroutines")
		c := make(chan int)
		quit := make(chan int)
		go func() {
			for i := 0; i < 10; i++ {
				fmt.Println(<-c)
				time.Sleep(100 * time.Millisecond)
			}
			quit <- 0
		}()
		fibonacci(c, quit)

		go func() {
			tick := time.Tick(100 * time.Millisecond)
			boom := time.After(500 * time.Millisecond)
			for {
				select {
				case <-tick:
					fmt.Println("tick.")
				case <-boom:
					fmt.Println("BOOM!")
					return
				default:
					fmt.Println("    .")
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()

		sc := SafeCounter{v: make(map[string]int)}
		for i := 0; i < 1000; i++ {
			go sc.Inc("somekey")
		}
		time.Sleep(time.Second)
		fmt.Println(sc.Value("somekey"))

	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		go func() {
			longOperation(ctx)
			cancel()
		}()

		<-ctx.Done()
		<-time.After(time.Second)
	}

}

func longOperation(ctx context.Context) {
	fmt.Println("starting long operation")
	defer func() {
		fmt.Println("exiting long operation")
	}()
	counter := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("operation halted")
			return
		default:
			counter++
			if counter > 3 {
				fmt.Println("long operation complete")
				return
			}
			time.Sleep(time.Second)
			fmt.Print(".")
		}
	}
}
