package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()
	fmt.Println(now.Format("03:04:05 PM"))
	centralTime := now.Add(-5 * time.Hour)
	fmt.Println(centralTime.Format("03:04:05 PM"))
}

// See if this will work
// Time subtraction
//
// Go allows you to subtract from times using the Add function with negative values. For example, if you want to deduct an hour from a given time, you can do that like so:
//
// package main
//
// import (
//     "fmt"
//     "time"
// )
//
// func main() {
//     givenDate := time.Now()
//     minusOneHour := givenDate.Add(-24 * time.Hour)
//     fmt.Println("One hour ago was: ", minusOneHour.Format(time.RFC822))
// }
//
// The code above subtracts 24 hours from the current time and prints a message based on the formatted result. The result should look like the following but based on the time at runtime:
//
// One hour ago was:  18 Jun 23 13:39 WAT
//
