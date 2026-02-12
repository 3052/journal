package main

import (
   "encoding/json"
   "flag"
   "log"
   "os"
   "path/filepath"
   "strings"
)

// FileData struct uses the clear "name" and "data" keys.
type FileData struct {
   Name string `json:"name"`
   Data string `json:"data"`
}

func main() {
   // 1. Define and parse the required -dir flag.
   inputDir := flag.String("dir", "", "The path to the input directory (required).")
   flag.Parse()

   if *inputDir == "" {
      log.Println("Error: The -dir flag is a required argument.")
      flag.Usage()
      os.Exit(1)
   }

   log.Printf("Target directory is: %s", *inputDir)

   // Validate that the provided path is a valid directory.
   info, err := os.Stat(*inputDir)
   if os.IsNotExist(err) {
      log.Fatalf("Error: The folder '%s' does not exist.", *inputDir)
   }
   if !info.IsDir() {
      log.Fatalf("Error: The path '%s' is not a directory.", *inputDir)
   }

   // 2. Find all processable files in that directory.
   sourceFiles, err := findSourceFiles(*inputDir)
   if err != nil {
      log.Fatalf("Error finding files in '%s': %v", *inputDir, err)
   }
   if len(sourceFiles) == 0 {
      log.Printf("Warning: No files were found in '%s'.", *inputDir)
      return
   }

   log.Printf("Found %d files in '%s' to process...", len(sourceFiles), *inputDir)

   // 3. Read the content of each file.
   var fileDataList []FileData
   for _, filename := range sourceFiles {
      fullPath := filepath.Join(*inputDir, filename)
      content, err := os.ReadFile(fullPath)
      if err != nil {
         log.Fatalf("Error reading file %s: %v", fullPath, err)
      }

      contentString := string(content)
      cleanedContent := strings.ReplaceAll(contentString, "\r", "")
      fileDataList = append(fileDataList, FileData{Name: filename, Data: cleanedContent})
   }

   // 4. Generate the compact JSON output.
   output, err := generateJSON(fileDataList)
   if err != nil {
      log.Fatalf("Error generating JSON output: %v", err)
   }

   // 5. --- THIS IS THE CORRECTED LINE ---
   // Write the output file to the CURRENT WORKING DIRECTORY, not the input directory.
   outputFilename := "combined.json"
   err = os.WriteFile(outputFilename, []byte(output), 0644)
   if err != nil {
      log.Fatalf("Error writing to file %s: %v", outputFilename, err)
   }

   log.Printf("Success! Output has been saved to %s", outputFilename)
}

// findSourceFiles now correctly uses os.ReadDir.
func findSourceFiles(targetDir string) ([]string, error) {
   entries, err := os.ReadDir(targetDir)
   if err != nil {
      return nil, err
   }

   var files []string
   for _, entry := range entries {
      if !entry.IsDir() {
         files = append(files, entry.Name())
      }
   }
   return files, nil
}

// generateJSON converts the file data into a compact, single-line JSON string.
func generateJSON(data []FileData) (string, error) {
   bytes, err := json.Marshal(data)
   if err != nil {
      return "", err
   }
   return string(bytes), nil
}
