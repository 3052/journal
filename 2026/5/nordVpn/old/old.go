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

func (c *client) do_country_code() error {
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
