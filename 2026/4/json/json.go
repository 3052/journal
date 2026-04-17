package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "os"
)

// prune traverses the JSON tree and truncates anything deeper than maxDepth.
func prune(v any, currentDepth, maxDepth int) any {
   // If we've exceeded the max depth, return a placeholder
   if currentDepth > maxDepth {
      return "[Truncated]"
   }

   switch val := v.(type) {
   case map[string]any:
      // It's a JSON object
      result := make(map[string]any)
      for k, v := range val {
         result[k] = prune(v, currentDepth+1, maxDepth)
      }
      return result
   case []any:
      // It's a JSON array
      result := make([]any, len(val))
      for i, v := range val {
         result[i] = prune(v, currentDepth+1, maxDepth)
      }
      return result
   default:
      // It's a primitive (string, number, bool, nil), return as is
      return val
   }
}

func main() {
   // Define the command-line flags
   inFile := flag.String("in", "", "Input JSON file path (required)")
   outFile := flag.String("out", "", "Output JSON file path (required)")
   maxDepth := flag.Int("depth", 2, "Maximum depth of JSON to output")
   
   flag.Parse()

   // Validate that required flags are provided
   if *inFile == "" || *outFile == "" {
      flag.Usage()
      os.Exit(1)
   }

   // Open the input file
   input, err := os.Open(*inFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
      os.Exit(1)
   }
   defer input.Close() // Ensure the file is closed when the function exits

   // Create the output file (this will overwrite the file if it already exists)
   output, err := os.Create(*outFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
      os.Exit(1)
   }
   defer output.Close() // Ensure the file is closed when the function exits

   // Use a Decoder to read from the input file
   decoder := json.NewDecoder(input)
   // UseNumber prevents large integers from being converted to scientific notation floats
   decoder.UseNumber()

   // Use an Encoder to write to the output file
   encoder := json.NewEncoder(output)
   encoder.SetIndent("", "  ") // Pretty-print the output

   // Loop to handle potential multiple JSON objects in the file (e.g., JSON Lines)
   for {
      var data any
      err := decoder.Decode(&data)
      if err == io.EOF {
         break
      }
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
         os.Exit(1)
      }

      // Prune the parsed JSON
      prunedData := prune(data, 0, *maxDepth)

      // Write the result to the output file
      if err := encoder.Encode(prunedData); err != nil {
         fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
         os.Exit(1)
      }
   }
   
   fmt.Printf("Successfully processed JSON to a max depth of %d and saved to %s\n", *maxDepth, *outFile)
}
