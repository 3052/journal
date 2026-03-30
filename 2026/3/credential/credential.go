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

func main() {
   // Define command-line flags
   host := flag.String("h", "", "Host to search for (e.g., amcplus.com)")
   key := flag.String("k", "", "Key to retrieve (e.g., password) - Optional")
   setFile := flag.String("set-file", "", "Save the JSON file location permanently")

   flag.Parse()

   // Execute core logic. If an error is returned, print to stderr and exit 1.
   if err := run(*host, *key, *setFile); err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      os.Exit(1)
   }
}

// run orchestrates the loading, validating, and searching of credentials
func run(host, key, setFile string) error {
   configPath, err := getConfigPath()
   if err != nil {
      return fmt.Errorf("locating user config directory: %w", err)
   }

   // 1. Handle saving the file location if --set-file is provided
   if setFile != "" {
      return saveConfig(setFile, configPath)
   }

   // 2. Load the data file location from the config
   dataFile, err := loadConfig(configPath)
   if err != nil {
      return err
   }
   if dataFile == "" {
      return fmt.Errorf("no data file location configured. Please use '--set-file <path>' first")
   }

   // 3. Validate host flag
   if host == "" {
      return fmt.Errorf("usage: credential -h <host> [-k <key>]")
   }

   // 4. Read the target credentials JSON file
   data, err := os.ReadFile(dataFile)
   if err != nil {
      return fmt.Errorf("reading credentials file '%s': %w", dataFile, err)
   }

   // 5. Parse the JSON
   var credentials []map[string]interface{}
   if err := json.Unmarshal(data, &credentials); err != nil {
      return fmt.Errorf("parsing JSON data: %w", err)
   }

   // 6. Validate the JSON data rules
   if err := validateData(credentials); err != nil {
      return err
   }

   // 7. Search and output
   return searchAndPrint(credentials, host, key)
}

// getConfigPath determines where to save/load the configuration file
func getConfigPath() (string, error) {
   configDir, err := os.UserConfigDir()
   if err != nil {
      return "", err
   }
   appConfigDir := filepath.Join(configDir, "credential")
   return filepath.Join(appConfigDir, "config.json"), nil
}

// saveConfig resolves the file path and saves it to the user's config directory
func saveConfig(setFile, configPath string) error {
   absPath, err := filepath.Abs(setFile)
   if err != nil {
      return fmt.Errorf("getting absolute path: %w", err)
   }

   cfg := AppConfig{DataFile: absPath}
   configData, err := json.MarshalIndent(cfg, "", "  ")
   if err != nil {
      return fmt.Errorf("encoding config data: %w", err)
   }

   if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
      return fmt.Errorf("creating config directory: %w", err)
   }

   if err := os.WriteFile(configPath, configData, 0644); err != nil {
      return fmt.Errorf("saving config file: %w", err)
   }

   // Print success to standard output
   fmt.Printf("Successfully saved data file location: %s\n", absPath)
   return nil
}

// loadConfig reads the config file and returns the saved DataFile path
func loadConfig(configPath string) (string, error) {
   b, err := os.ReadFile(configPath)
   if err != nil {
      if os.IsNotExist(err) {
         return "", nil // Normal if the user hasn't run --set-file yet
      }
      return "", fmt.Errorf("reading config file: %w", err)
   }

   var cfg AppConfig
   if err := json.Unmarshal(b, &cfg); err != nil {
      return "", fmt.Errorf("parsing config file (it might be corrupted): %w", err)
   }

   return cfg.DataFile, nil
}

// searchAndPrint finds the matching object(s) and prints them to standard output
func searchAndPrint(credentials []map[string]interface{}, host, key string) error {
   if key != "" {
      // Specific Key requested: Print just the value
      for _, cred := range credentials {
         if credHost, ok := cred["host"].(string); ok && credHost == host {
            if val, exists := cred[key]; exists {
               fmt.Printf("%v\n", val)
               return nil
            }
         }
      }
      return fmt.Errorf("could not find key '%s' for host '%s'", key, host)
   }

   // No Key requested: Collect all objects matching the host
   var matches []map[string]interface{}
   for _, cred := range credentials {
      if credHost, ok := cred["host"].(string); ok && credHost == host {
         matches = append(matches, cred)
      }
   }

   if len(matches) == 0 {
      return fmt.Errorf("could not find any entries for host '%s'", host)
   }

   // Format and print the matching objects as JSON
   var out []byte
   var err error

   if len(matches) == 1 {
      out, err = json.MarshalIndent(matches[0], "", "  ")
   } else {
      out, err = json.MarshalIndent(matches, "", "  ")
   }

   if err != nil {
      return fmt.Errorf("formatting output JSON: %w", err)
   }

   fmt.Println(string(out))
   return nil
}

// validateData enforces the specific data rules on the JSON credentials
func validateData(credentials []map[string]interface{}) error {
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
            isTrialFalse = !v
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
