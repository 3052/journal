package openRouter

import (
   "io"
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := provider_filters()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("provider_filters.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
