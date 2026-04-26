package main

import (
   "bytes"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "os"
)

func main() {
   // 1. Define command-line flags
   inputFile := flag.String("i", "", "Path to the input JSON file (required)")
   outputFile := flag.String("o", "", "Path to the output JSON file (required)")
   // 2. Parse the flags
   flag.Parse()
   // 3. Ensure both flags were provided
   if *inputFile == "" || *outputFile == "" {
      flag.Usage()
      os.Exit(1)
   }
   // 4. Read the input file
   inputJSON, err := os.ReadFile(*inputFile)
   if err != nil {
      log.Fatalf("Failed to read input file: %v", err)
   }
   // 5. Create a buffer and compact the JSON
   var compacted bytes.Buffer
   err = json.Compact(&compacted, inputJSON)
   if err != nil {
      log.Fatalf("Invalid JSON input: %v", err)
   }
   // 6. Write the compacted JSON to the output file using os.ModePerm
   err = os.WriteFile(*outputFile, compacted.Bytes(), os.ModePerm)
   if err != nil {
      log.Fatalf("Failed to write to output file: %v", err)
   }
   fmt.Printf("Successfully compacted '%s' into '%s'\n", *inputFile, *outputFile)
}
