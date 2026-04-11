package main

import (
   "flag"
   "fmt"
   "os"
   "path/filepath"
   "strings"
)

func main() {
   // Define command-line flags. Default dir is empty to prevent accidental execution.
   dirPtr := flag.String("dir", "", "The directory to process (REQUIRED)")
   restorePtr := flag.Bool("restore", false, "Set to true to REMOVE the .txt extensions (undo)")
   flag.Parse()

   // Safety check: Do nothing if no directory is explicitly provided
   if *dirPtr == "" {
      fmt.Println("Error: You must provide a directory.")
      fmt.Println("Usage: go run rename.go -dir <path_to_folder>")
      os.Exit(1)
   }

   // Walk through the directory recursively
   err := filepath.WalkDir(*dirPtr, func(path string, d os.DirEntry, err error) error {
      if err != nil {
         return err
      }

      // Skip directories completely
      if d.IsDir() {
         // Safety feature: skip .git folders so we don't break repositories
         if d.Name() == ".git" {
            return filepath.SkipDir
         }
         return nil
      }

      // RESTORE MODE: Remove ".txt"
      if *restorePtr {
         if strings.HasSuffix(d.Name(), ".txt") {
            newPath := strings.TrimSuffix(path, ".txt")
            fmt.Printf("Restoring: %s  ->  %s\n", filepath.Base(path), filepath.Base(newPath))
            return os.Rename(path, newPath)
         }
         return nil
      }

      // RENAME MODE: Add ".txt"
      // Skip files that already end in .txt to prevent file.go.txt.txt
      if !strings.HasSuffix(d.Name(), ".txt") {
         newPath := path + ".txt"
         fmt.Printf("Renaming: %s  ->  %s\n", filepath.Base(path), filepath.Base(newPath))
         return os.Rename(path, newPath)
      }

      return nil
   })

   if err != nil {
      fmt.Printf("Error processing directory: %v\n", err)
   } else {
      fmt.Println("\nDone!")
   }
}
