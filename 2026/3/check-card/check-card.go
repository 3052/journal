package main

import (
   "flag"
   "log"
   "os"
   "strings"
)

func read_file(name string) (string, error) {
   file, err := os.Open(name)
   if err != nil {
      return "", err
   }
   defer file.Close()
   var data strings.Builder
   _, err = file.WriteTo(&data)
   if err != nil {
      return "", err
   }
   return data.String(), nil
}

var contains = []string{
   "BRAUMS STORE",
   "CHICK-FIL-A",
   "JASON'S DELI",
   "LA MADELEINE",
   "MCDONALD'S",
   "SCHLOTZSKYS",
   "SPRING CREEK",
   "WENDY",
   "WHATABURGER",
}

func do_check(name string) error {
   data, err := read_file(name)
   if err != nil {
      return err
   }
   for _, contain := range contains {
      if strings.Contains(data, contain) {
         log.Println("Contains", contain)
      }
   }
   log.Print()
   for _, contain := range contains {
      if !strings.Contains(data, contain) {
         log.Println("!Contains", contain)
      }
   }
   return nil
}

func main() {
   log.SetFlags(log.Ltime)
   name := flag.String("n", "", "name")
   flag.Parse()
   if *name != "" {
      err := do_check(*name)
      if err != nil {
         log.Fatal(err)
      }
   } else {
      flag.Usage()
   }
}
