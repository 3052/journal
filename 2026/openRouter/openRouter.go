package openRouter

import (
   "net/http"
)

func provider_filters() (*http.Response, error) {
   return http.Get("https://openrouter.ai/api/frontend/provider-filters")
}
