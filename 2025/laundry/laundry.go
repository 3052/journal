package main

import "fmt"

var basket = []string{
   "outside",
   "inside floor",
   "inside shelf",
}

var door = []string{
   "locked",
   "unlocked",
}

func main() {
   for _, a := range basket {
      for _, b := range door {
         fmt.Println(a, b)
      }
   }
}
