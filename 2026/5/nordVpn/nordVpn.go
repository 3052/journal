package main

import (
   "encoding/json"
   "fmt"
   "sort"
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
func GetFastestServers(servers []Server, countryCode string, max int) []Server {
   var filtered []Server

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

   // 2. Sort the filtered servers by Load (ascending = fastest/least busy)
   sort.Slice(filtered, func(i, j int) bool {
      return filtered[i].Load < filtered[j].Load
   })

   // 3. Limit to the max requested servers (e.g., 9)
   if len(filtered) > max {
      return filtered[:max]
   }
   return filtered
}

func main() {
   // Sample JSON data based on your provided file
   jsonData := `[
      {
         "id": 929966,
         "name": "Poland #128",
         "hostname": "pl128.nordvpn.com",
         "load": 40,
         "status": "online",
         "locations":[
            {"country": {"code": "PL"}}
         ]
      },
      {
         "id": 929967,
         "name": "Poland #129",
         "hostname": "pl129.nordvpn.com",
         "load": 12,
         "status": "online",
         "locations": [
            {"country": {"code": "PL"}}
         ]
      },
      {
         "id": 929968,
         "name": "Germany #10",
         "hostname": "de10.nordvpn.com",
         "load": 5,
         "status": "online",
         "locations":[
            {"country": {"code": "DE"}}
         ]
      }
   ]`

   // Unmarshal JSON into a slice of Server structs
   var servers []Server
   if err := json.Unmarshal([]byte(jsonData), &servers); err != nil {
      panic(err)
   }

   // Request fastest servers for "PL", max 9
   targetCountry := "PL"
   maxServers := 9

   fastest := GetFastestServers(servers, targetCountry, maxServers)

   // Print the results
   fmt.Printf("Top %d fastest servers for %s:\n", maxServers, targetCountry)
   fmt.Println("-------------------------------------------------")
   for i, s := range fastest {
      fmt.Printf("%d. %s (%s) - Load: %d%%\n", i+1, s.Name, s.Hostname, s.Load)
   }
}
