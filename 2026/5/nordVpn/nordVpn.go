package main

import (
   "cmp"
   "encoding/json"
   "fmt"
   "log"
   "os"
   "path/filepath"
   "slices"
   "time"
)

// Define the necessary structs to parse the JSON data
type Country struct {
   Code string `json:"code"`
}

type Location struct {
   Country Country `json:"country"`
}

type Server struct {
   ID        int        `json:"id"`
   Name      string     `json:"name"`
   Hostname  string     `json:"hostname"`
   Load      int        `json:"load"`
   Status    string     `json:"status"`
   Locations []Location `json:"locations"`
}

// GetFastestServers filters by country code, sorts by lowest load, and limits the result.
func GetFastestServers(servers []*Server, countryCode string, limit int) []*Server {
   var filtered []*Server

   // 1. Filter by country code and online status
   for _, s := range servers {
      if s.Status == "online" {
         for _, loc := range s.Locations {
            if loc.Country.Code == countryCode {
               filtered = append(filtered, s)
               break
            }
         }
      }
   }

   // 2. Sort the filtered pointers by Load
   slices.SortFunc(filtered, func(a, b *Server) int {
      return cmp.Compare(a.Load, b.Load)
   })

   // 3. Limit the results
   if len(filtered) > limit {
      return filtered[:limit]
   }
   return filtered
}

func main() {
   // 1. Get the user's cache directory
   cacheDir, err := os.UserCacheDir()
   if err != nil {
      log.Fatalf("Failed to get user cache directory: %v", err)
   }

   // 2. Safely construct the file path
   filePath := filepath.Join(cacheDir, "nordVpn", "nordVpn.json")

   // 3. Check if the file exists and how old it is
   fileInfo, err := os.Stat(filePath)
   if err != nil {
      log.Fatalf("Failed to access file %s: %v", filePath, err)
   }

   // If the file was modified 24 hours ago or more, return an error
   if time.Since(fileInfo.ModTime()) >= 24*time.Hour {
      log.Fatalf("Error: the file %s is 24 hours old or more. Please update it.", filePath)
   }

   // 4. Read the JSON file
   jsonData, err := os.ReadFile(filePath)
   if err != nil {
      log.Fatalf("Failed to read file %s: %v", filePath, err)
   }

   // 5. Unmarshal JSON directly into a slice of pointers
   var servers []*Server
   if err := json.Unmarshal(jsonData, &servers); err != nil {
      log.Fatalf("Failed to parse JSON: %v", err)
   }

   // 6. Request fastest servers for "PL", max limit 9
   targetCountry := "PL"
   limit := 9

   fastest := GetFastestServers(servers, targetCountry, limit)

   // 7. Print the results
   if len(fastest) == 0 {
      fmt.Printf("No online servers found for %s.\n", targetCountry)
      return
   }

   fmt.Printf("Top %d fastest servers for %s:\n", limit, targetCountry)
   fmt.Println("-------------------------------------------------")
   for i, s := range fastest {
      fmt.Printf("%d. %s (%s) - Load: %d%%\n", i+1, s.Name, s.Hostname, s.Load)
   }
}
