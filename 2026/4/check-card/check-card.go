package main

import (
   "bufio"
   "cmp"
   "flag"
   "fmt"
   "os"
   "regexp"
   "slices"
   "strings"
   "time"
)

// MatchResult now holds a Count instead of a boolean
type MatchResult struct {
   Description string
   Count       int
}

func main() {
   // Define command-line flag for the input file
   inputFile := flag.String("f", "", "Path to the bank statement text file")
   flag.Parse()
   // Ensure the file flag was provided
   if *inputFile == "" {
      flag.Usage()
      os.Exit(1)
   }

   // 1. Define the target keywords to check against
   targets := []string{
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

   // 2. Set timeframe boundaries normalized to midnight
   now := time.Now()
   today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
   oneWeekAgo := today.AddDate(0, 0, -7)

   // 3. Open the bank statement file using the provided flag
   file, err := os.Open(*inputFile)
   if err != nil {
      panic(err)
   }
   defer file.Close()

   var validDescriptions []string
   scanner := bufio.NewScanner(file)

   // 4. Regex strictly for MM/DD format
   dateRegex := regexp.MustCompile(`\b(\d{2}/\d{2})\b`)

   // Parse file to extract matching descriptions
   for scanner.Scan() {
      line := strings.TrimSpace(scanner.Text())

      if line == "Description" {
         if scanner.Scan() {
            desc := strings.TrimSpace(scanner.Text())

            // Extract the MM/DD date from the description string
            match := dateRegex.FindString(desc)
            if match != "" {
               // Parse just the MM/DD to extract the month and day
               parsedMMDD, err := time.Parse("01/02", match)

               if err == nil {
                  year := today.Year()

                  // If the transaction month is greater than the current month,
                  // the week crossed January 1st (e.g., Today is Jan, transaction is Dec).
                  if parsedMMDD.Month() > today.Month() {
                     year--
                  }

                  // Create the actual date object using the evaluated year
                  parsedDate := time.Date(year, parsedMMDD.Month(), parsedMMDD.Day(), 0, 0, 0, 0, today.Location())

                  // Filter: Description date must be strictly on or after the 7-day cutoff.
                  if !parsedDate.Before(oneWeekAgo) {
                     // Append directly without strings.ToUpper
                     validDescriptions = append(validDescriptions, desc)
                  }
               }
            }
         }
      }
   }

   if err := scanner.Err(); err != nil {
      panic(err)
   }

   // 5. Evaluate which targets are found in the filtered records and count them
   var results []*MatchResult
   for _, target := range targets {
      count := 0
      for _, desc := range validDescriptions {
         if strings.Contains(desc, target) {
            count++
         }
      }
      results = append(results, &MatchResult{
         Description: target,
         Count:       count,
      })
   }

   // 6. Sort by Count descending (highest count first)
   // By comparing b to a (instead of a to b), we reverse the default ascending order.
   slices.SortFunc(results, func(a, b *MatchResult) int {
      return cmp.Compare(a.Count, b.Count)
   })

   // 7. Print the formatted Markdown table
   fmt.Println("| count | description |")
   fmt.Println("|---|---|")
   for _, res := range results {
      fmt.Printf("| %d | %s |\n", res.Count, res.Description)
   }
}
