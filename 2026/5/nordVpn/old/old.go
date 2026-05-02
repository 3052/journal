package main

import (
   "encoding/json"
   "errors"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "path/filepath"
   "strconv"
   "strings"
   "time"
)

// limit <= -1 for default
// limit == 0 for all
func WriteServers(limit int) ([]byte, error) {
   var query string
   if limit >= 0 {
      query = "limit=" + strconv.Itoa(limit)
   }
   resp, err := Get(
      &url.URL{
         Scheme:   "https",
         Host:     "api.nordvpn.com",
         Path:     "/v1/servers",
         RawQuery: query,
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func FormatProxy(username, password, hostname string) string {
   var data strings.Builder
   data.WriteString("https://")
   data.WriteString(username)
   data.WriteByte(':')
   data.WriteString(password)
   data.WriteByte('@')
   data.WriteString(hostname)
   data.WriteString(":89")
   return data.String()
}

func ReadServers(data []byte) ([]Server, error) {
   var result []Server
   err := json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (s *Server) ProxySsl() bool {
   for _, technology := range s.Technologies {
      if technology.Identifier == "proxy_ssl" {
         return true
      }
   }
   return false
}

func (s *Server) Country(code string) bool {
   for _, location := range s.Locations {
      if location.Country.Code == code {
         return true
      }
   }
   return false
}

type Server struct {
   Hostname     string
   Status       string
   Technologies []struct {
      Identifier string
   }
   Locations []struct {
      Country struct {
         City struct {
            DnsName string `json:"dns_name"`
         }
         Code string
      }
   }
}

func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error) {
   reqHeader := make(http.Header)
   for key, value := range headers {
      reqHeader.Set(key, value)
   }
   req := &http.Request{
      Method: http.MethodGet,
      URL:    targetUrl,
      Header: reqHeader,
   }

   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   var err error
   c.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   c.cache = filepath.Join(c.cache, "nordVpn/nordVpn.json")
   // 1
   flag.BoolVar(&c.write, "w", false, "write")
   // 2
   flag.StringVar(&c.country_code, "c", "", "country code")
   flag.Parse()
   if c.write {
      return c.do_write()
   }
   if c.country_code != "" {
      return c.do_country_code()
   }
   flag.Usage()
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *client) do_write() error {
   data, err := WriteServers(0)
   if err != nil {
      return err
   }
   return write_file(c.cache, data)
}

type client struct {
   cache string
   // 1
   write bool
   // 2
   country_code string
}

func (c *client) do_country_code() error {
   data, err := read_file(c.cache)
   if err != nil {
      return err
   }
   servers, err := ReadServers(data)
   if err != nil {
      return err
   }
   data, err = exec.Command("credential", "-j=api.nordvpn.com").Output()
   if err != nil {
      return err
   }
   var credential []struct {
      Username string
      Password string
   }
   err = json.Unmarshal(data, &credential)
   if err != nil {
      return err
   }
   for _, server := range servers {
      if server.ProxySsl() {
         if server.Country(c.country_code) {
            fmt.Println(
               FormatProxy(
                  credential[0].Username, credential[0].Password,
                  server.Hostname,
               ),
            )
         }
      }
   }
   return nil
}

const duration = 24 * time.Hour

func read_file(name string) ([]byte, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   info, err := file.Stat()
   if err != nil {
      return nil, err
   }
   if time.Since(info.ModTime()) >= duration {
      return nil, errors.New(duration.String())
   }
   return io.ReadAll(file)
}
