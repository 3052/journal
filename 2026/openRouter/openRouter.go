package openRouter

import (
   "encoding/json"
   "net/http"
   "net/url"
)

func (p provider_filter) find() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "openrouter.ai",
      Path: "/api/frontend/models/find",
      RawQuery: url.Values{
         //context=16000&
         //fmt=cards&
         //input_modalities=text&
         //output_modalities=text&
         "providers": {p.Name},
      }.Encode(),
   }
   return http.DefaultClient.Do(&req)
}

type provider_filter struct {
   Name string
}

func provider_filters() ([]provider_filter, error) {
   resp, err := http.Get("https://openrouter.ai/api/frontend/provider-filters")
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data []provider_filter
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Data, nil
}
