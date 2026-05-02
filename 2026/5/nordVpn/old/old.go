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
