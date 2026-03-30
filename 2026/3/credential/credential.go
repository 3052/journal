package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "os"
   "path/filepath"
   "time"
)

// AppConfig stores the user's saved preferences
type AppConfig struct {
   DataFile string `json:"data_file"`
}

// getConfigPath determines where to save/load the configuration file
func getConfigPath() (string, error) {
   configDir, err := os.UserConfigDir()
   if err != nil {
      return "", fmt.Errorf("getting user config dir: %w", err)
   }
   appConfigDir := filepath.Join(configDir, "credential")
   return filepath.Join(appConfigDir, "config.json"), nil
}

// validateData enforces the specific data rules on the JSON credentials
func validateData(credentials []map[string]interface{}) error {
   // Map to count how many times each password is used across all objects
   passCounts := make(map[string]int)

   // First pass: Count password occurrences
   for _, cred := range credentials {
      if passVal, exists := cred["password"]; exists {
         if passStr, ok := passVal.(string); ok {
            passCounts[passStr]++
         }
      }
   }

   // Calculate the date exactly 1 year ago from today
   oneYearAgo := time.Now().AddDate(-1, 0, 0)

   // Second pass: Validate each rule
   for i, cred := range credentials {
      hostStr, ok := cred["host"].(string)
      if !ok {
         hostStr = "unknown/missing host"
      }

      // --- Rule 3: All objects must have a date, and it cannot be older than 1 year ---
      dateVal, dateExists := cred["date"]
      if !dateExists {
         return fmt.Errorf("validation error: object at index %d (host: %s) is missing 'date' key", i, hostStr)
      }

      dateStr, ok := dateVal.(string)
      if !ok {
         return fmt.Errorf("validation error: object at index %d (host: %s) has invalid 'date' format", i, hostStr)
      }

      // Parse date assuming YYYY-MM-DD format
      parsedDate, err := time.Parse("2006-01-02", dateStr)
      if err != nil {
         return fmt.Errorf("validation error: object at index %d (host: %s) has invalid date '%s' (expected YYYY-MM-DD)", i, hostStr, dateStr)
      }
      if parsedDate.Before(oneYearAgo) {
         return fmt.Errorf("validation error: object at index %d (host: %s) has a date '%s' older than 1 year", i, hostStr, dateStr)
      }

      // --- Rule 1: If object has "password", it must have "trial" ---
      passVal, passExists := cred["password"]
      if passExists {
         if _, trialExists := cred["trial"]; !trialExists {
            return fmt.Errorf("validation error: object at index %d (host: %s) has a 'password' but is missing the 'trial' key", i, hostStr)
         }
      }

      // --- Rule 2: For trial=false objects, password cannot match any other objects ---
      if trialVal, trialExists := cred["trial"]; trialExists {
         isTrialFalse := false

         // Handle both string "false" and boolean false JSON values
         switch v := trialVal.(type) {
         case string:
            isTrialFalse = (v == "false")
         case bool:
            isTrialFalse = (v == false)
         }

         if isTrialFalse && passExists {
            if passStr, ok := passVal.(string); ok {
               if passCounts[passStr] > 1 {
                  return fmt.Errorf("validation error: trial=false object at index %d (host: %s) shares its password with another object", i, hostStr)
               }
            }
         }
      }
   }

   return nil
}

func main() {
   // Define command-line flags
   host := flag.String("h", "", "Host to search for (e.g., amcplus.com)")
   key := flag.String("k", "", "Key to retrieve (e.g., password) - Optional")
   setFile := flag.String("set-file", "", "Save the JSON file location permanently")

   flag.Parse()

   // 1. Get the path to the user's config file
   configPath, err := getConfigPath()
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error locating user config directory: %v\n", err)
      os.Exit(1)
   }

   // 2. Handle saving the file location if --set-file is provided
   if *setFile != "" {
      absPath, err := filepath.Abs(*setFile)
      if err != nil {
         fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
         os.Exit(1)
      }

      cfg := AppConfig{DataFile: absPath}
      configData, err := json.MarshalIndent(cfg, "", "  ")
      if err != nil {
         fmt.Fprintf(os.Stderr, "Failed to encode config data: %v\n", err)
         os.Exit(1)
      }

      if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
         fmt.Fprintf(os.Stderr, "Failed to create config directory: %v\n", err)
         os.Exit(1)
      }

      if err := os.WriteFile(configPath, configData, 0644); err != nil {
         fmt.Fprintf(os.Stderr, "Failed to save config file: %v\n", err)
         os.Exit(1)
      }

      fmt.Printf("Successfully saved data file location: %s\n", absPath)
      os.Exit(0)
   }

   // 3. Load the data file location from the config
   var dataFile string
   b, err := os.ReadFile(configPath)
   if err != nil {
      // It's okay if the config file doesn't exist yet, but other errors should be reported
      if !os.IsNotExist(err) {
         fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
         os.Exit(1)
      }
   } else {
      var cfg AppConfig
      if err := json.Unmarshal(b, &cfg); err != nil {
         fmt.Fprintf(os.Stderr, "Error parsing config file (it might be corrupted): %v\n", err)
         os.Exit(1)
      }
      dataFile = cfg.DataFile
   }

   if dataFile == "" {
      fmt.Fprintln(os.Stderr, "Error: No data file location configured.")
      fmt.Fprintln(os.Stderr, "Please use '--set-file <path>' to save your JSON location first.")
      os.Exit(1)
   }

   // 4. Validate host flag
   if *host == "" {
      fmt.Fprintln(os.Stderr, "Usage: credential -h <host> [-k <key>]")
      os.Exit(1)
   }

   // 5. Read the target credentials JSON file
   data, err := os.ReadFile(dataFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading credentials file '%s': %v\n", dataFile, err)
      os.Exit(1)
   }

   // 6. Parse the JSON
   var credentials []map[string]interface{}
   if err := json.Unmarshal(data, &credentials); err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing JSON data: %v\n", err)
      os.Exit(1)
   }

   // 7. Validate the JSON data rules
   if err := validateData(credentials); err != nil {
      fmt.Fprintf(os.Stderr, "%v\n", err)
      os.Exit(1)
   }

   // 8. Search Logic
   if *key != "" {
      // Specific Key requested: Print just the value
      for _, cred := range credentials {
         if credHost, ok := cred["host"].(string); ok && credHost == *host {
            if val, exists := cred[*key]; exists {
               fmt.Printf("%v\n", val)
               os.Exit(0)
            }
         }
      }
      fmt.Fprintf(os.Stderr, "Could not find key '%s' for host '%s'\n", *key, *host)
      os.Exit(1)

   } else {
      // No Key requested: Collect all objects matching the host
      var matches []map[string]interface{}
      for _, cred := range credentials {
         if credHost, ok := cred["host"].(string); ok && credHost == *host {
            matches = append(matches, cred)
         }
      }

      if len(matches) == 0 {
         fmt.Fprintf(os.Stderr, "Could not find any entries for host '%s'\n", *host)
         os.Exit(1)
      }

      // Print the matching objects as formatted JSON
      var out []byte
      var marshalErr error

      if len(matches) == 1 {
         out, marshalErr = json.MarshalIndent(matches[0], "", "  ")
      } else {
         out, marshalErr = json.MarshalIndent(matches, "", "  ")
      }

      if marshalErr != nil {
         fmt.Fprintf(os.Stderr, "Failed to format output JSON: %v\n", marshalErr)
         os.Exit(1)
      }

      fmt.Println(string(out))
      os.Exit(0)
   }
}
