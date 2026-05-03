package main

import (
   "cmp"
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "slices"
   "strings"
   "time"
)

// Define the necessary structs to parse the JSON data
type Country struct {
   Code string `json:"code"`
}

type Location struct {
   Country Country `json:"country"`
}

type Technology struct {
   ID         int    `json:"id"`
   Name       string `json:"name"`
   Identifier string `json:"identifier"`
}

type Server struct {
   ID           int          `json:"id"`
   Name         string       `json:"name"`
   Hostname     string       `json:"hostname"`
   Load         int          `json:"load"`
   Status       string       `json:"status"`
   Locations    []Location   `json:"locations"`
   Technologies []Technology `json:"technologies"`
}

// GetFastestServers filters by country code, checks for proxy_ssl, sorts by lowest load, and limits the result.
func GetFastestServers(servers []*Server, countryCode string, limit int) []*Server {
   var filtered []*Server

   // 1. Filter by country code, online status, and proxy_ssl technology
   for _, s := range servers {
      if s.Status != "online" {
         continue
      }

      // Check if the server is in the target country
      hasCountry := false
      for _, loc := range s.Locations {
         if strings.EqualFold(loc.Country.Code, countryCode) {
            hasCountry = true
            break
         }
      }
      if !hasCountry {
         continue
      }

      // Check if the server supports proxy_ssl
      hasProxySSL := false
      for _, tech := range s.Technologies {
         if tech.Identifier == "proxy_ssl" {
            hasProxySSL = true
            break
         }
      }
      if !hasProxySSL {
         continue
      }

      // If it passes all checks, add it to our filtered list
      filtered = append(filtered, s)
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

func refreshFile(filePath string) error {
   u := url.URL{
      Scheme:   "https",
      Host:     "api.nordvpn.com",
      Path:     "/v1/servers",
      RawQuery: "limit=0",
   }

   fmt.Printf("Downloading latest server list from %s...\n", u.String())
   resp, err := http.Get(u.String())
   if err != nil {
      return fmt.Errorf("failed to fetch data: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
   }

   // Ensure the cache directory exists before writing
   if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
      return fmt.Errorf("failed to create directory: %w", err)
   }

   out, err := os.Create(filePath)
   if err != nil {
      return fmt.Errorf("failed to create file: %w", err)
   }
   defer out.Close()

   if _, err := io.Copy(out, resp.Body); err != nil {
      return fmt.Errorf("failed to write data to file: %w", err)
   }

   fmt.Println("Server list successfully updated.")
   return nil
}

func main() {
   // Setup the command line flags
   refresh := flag.Bool("refresh", false, "Fetch the latest server list from NordVPN")
   country := flag.String("country", "", "Target country code (e.g., PL, DE, US) [REQUIRED]")
   flag.Parse()

   // Enforce the requirement of the country flag
   if *country == "" {
      fmt.Fprintf(os.Stderr, "Error: the -country flag is required.\n\n")
      flag.Usage()
      os.Exit(1)
   }

   // 1. Get the user's cache directory
   cacheDir, err := os.UserCacheDir()
   if err != nil {
      log.Fatalf("Failed to get user cache directory: %v", err)
   }

   // 2. Safely construct the file path
   filePath := filepath.Join(cacheDir, "nordVpn", "nordVpn.json")

   // 3. Handle the refresh flag or file age check
   if *refresh {
      if err := refreshFile(filePath); err != nil {
         log.Fatalf("Refresh failed: %v", err)
      }
   } else {
      fileInfo, err := os.Stat(filePath)
      if err != nil {
         if os.IsNotExist(err) {
            log.Fatalf("Cache file does not exist. Please run the program with the -refresh flag.")
         }
         log.Fatalf("Failed to access file %s: %v", filePath, err)
      }

      // If the file was modified 24 hours ago or more, prompt user to refresh
      if time.Since(fileInfo.ModTime()) >= 24*time.Hour {
         log.Fatalf("Error: the file %s is 24 hours old or more. Please run with the -refresh flag.", filePath)
      }
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

   // 6. Request fastest servers, max limit 9
   limit := 9
   targetCountry := strings.ToUpper(*country)

   fastest := GetFastestServers(servers, targetCountry, limit)

   // 7. Print the results
   if len(fastest) == 0 {
      fmt.Printf("No online 'proxy_ssl' servers found for %s.\n", targetCountry)
      return
   }

   fmt.Printf("Top %d fastest 'proxy_ssl' servers for %s:\n", limit, targetCountry)
   fmt.Println("-------------------------------------------------")
   for i, s := range fastest {
      fmt.Printf("%d. %s (%s) - Load: %d%%\n", i+1, s.Name, s.Hostname, s.Load)
   }
}
