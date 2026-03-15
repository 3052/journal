package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "regexp"
   "strings"
   "sync"
)

var (
   // Authentication token optionally populated via environment variable
   gitToken string

   // Concurrency management
   wg  sync.WaitGroup
   // Limit to 9 concurrent file downloads at a time to prevent rate-limit / socket exhaustion
   sem = make(chan struct{}, 9)
)

// GitHubContent represents the JSON object returned by the GitHub Contents API
type GitHubContent struct {
   Name string `json:"name"`
   Path string `json:"path"`
   Type string `json:"type"`
   URL  string `json:"url"`
}

func main() {
   repoURL := flag.String("u", "", "The URL of the repository (e.g., https://github.com/owner/repo)")
   dirPath := flag.String("d", "", "The directory path inside the repo to download")
   flag.Parse()

   if *repoURL == "" || *dirPath == "" {
      fmt.Println("Error: Both -u (URL) and -d (directory) flags are required.")
      flag.Usage()
      os.Exit(1)
   }

   gitToken = os.Getenv("GIT_TOKEN")

   domain, repo, err := parseRepoURL(*repoURL)
   if err != nil {
      fmt.Printf("Error parsing URL: %v\n", err)
      os.Exit(1)
   }

   fmt.Printf("Starting concurrent download for directory '%s' from repo '%s'...\n", *dirPath, repo)

   // Begin recursive directory traversal
   err = downloadDir(domain, repo, *dirPath)
   if err != nil {
      fmt.Printf("Error processing directory: %v\n", err)
   }

   // Wait for all concurrent file downloads to finish
   wg.Wait()
   fmt.Println("Download complete.")
}

// parseRepoURL extracts the domain and the repository identifier (owner/repo) from a git URL
func parseRepoURL(rawURL string) (string, string, error) {
   // Handle SSH URLs like git@github.com:owner/repo.git
   if strings.HasPrefix(rawURL, "git@") {
      re := regexp.MustCompile(`^git@([^:]+):(.+?)(?:\.git)?$`)
      matches := re.FindStringSubmatch(rawURL)
      if len(matches) == 3 {
         return matches[1], matches[2], nil
      }
      return "", "", fmt.Errorf("invalid SSH URL format")
   }

   // Handle HTTP/HTTPS URLs
   u, err := url.Parse(rawURL)
   if err != nil {
      return "", "", err
   }
   domain := u.Host
   repo := strings.TrimPrefix(u.Path, "/")
   repo = strings.TrimSuffix(repo, ".git")
   return domain, repo, nil
}

// getAPIURL formats the GitHub API URL for standard GitHub or GitHub Enterprise
func getAPIURL(domain, repo, path string) string {
   if strings.ToLower(domain) == "github.com" {
      return fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
   }
   return fmt.Sprintf("https://%s/api/v3/repos/%s/contents/%s", domain, repo, path)
}

// doRequest is a wrapper to inject headers required by GitHub's API
func doRequest(reqURL, acceptHeader string) (*http.Response, error) {
   req, err := http.NewRequest("GET", reqURL, nil)
   if err != nil {
      return nil, err
   }
   if gitToken != "" {
      req.Header.Set("Authorization", "token "+gitToken)
   }
   if acceptHeader != "" {
      req.Header.Set("Accept", acceptHeader)
   } else {
      req.Header.Set("Accept", "application/vnd.github.v3+json")
   }
   return http.DefaultClient.Do(req)
}

// downloadDir handles fetching directory listings and queuing files or subdirectories
func downloadDir(domain, repo, dirPath string) error {
   apiURL := getAPIURL(domain, repo, dirPath)

   resp, err := doRequest(apiURL, "")
   if err != nil {
      return fmt.Errorf("request failed: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      bodyBytes, _ := io.ReadAll(resp.Body)
      return fmt.Errorf("API returned status %s: %s", resp.Status, string(bodyBytes))
   }

   var contents []GitHubContent
   if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
      return fmt.Errorf("failed to decode JSON response: %v", err)
   }

   for _, item := range contents {
      if item.Type == "dir" {
         // Traverse subdirectories synchronously to map out the tree without overloading limits
         if err := downloadDir(domain, repo, item.Path); err != nil {
            fmt.Printf("Error traversing directory %s: %v\n", item.Path, err)
         }
      } else if item.Type == "file" {
         // Schedule file downloading concurrently
         wg.Add(1)
         go func(fileURL, filePath string) {
            defer wg.Done()
            
            // Acquire semaphore
            sem <- struct{}{}
            defer func() { <-sem }()

            if err := downloadFile(fileURL, filePath); err != nil {
               fmt.Printf("[ERROR] Downloading %s: %v\n", filePath, err)
            } else {
               fmt.Printf("[OK] Downloaded %s\n", filePath)
            }
         }(item.URL, item.Path)
      }
   }
   return nil
}

// downloadFile writes the raw file payload directly to disk
func downloadFile(fileURL, localPath string) error {
   // Ensure the parent directories exist on the local file system
   dir := filepath.Dir(localPath)
   if err := os.MkdirAll(dir, os.ModePerm); err != nil {
      return fmt.Errorf("failed to create local directory: %v", err)
   }

   // Passing application/vnd.github.v3.raw returns the raw file bytes directly
   resp, err := doRequest(fileURL, "application/vnd.github.v3.raw")
   if err != nil {
      return fmt.Errorf("failed to make request: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("API returned status %s", resp.Status)
   }

   out, err := os.Create(localPath)
   if err != nil {
      return fmt.Errorf("failed to create local file: %v", err)
   }
   defer out.Close()

   _, err = io.Copy(out, resp.Body)
   if err != nil {
      return fmt.Errorf("failed to write content to disk: %v", err)
   }

   return nil
}
