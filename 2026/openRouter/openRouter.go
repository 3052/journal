package openRouter

import (
   "encoding/json"
   "net/http"
   "net/url"
)

type model_data struct {
   Name string
   Slug string
}

func (p provider_filter) find() ([]*model_data, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "openrouter.ai",
      Path: "/api/frontend/models/find",
      RawQuery: url.Values{"providers": {p.Name}}.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Models []*model_data
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Data.Models, nil
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
