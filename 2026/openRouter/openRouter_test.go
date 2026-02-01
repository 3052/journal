package openRouter

import (
   "io"
   "os"
   "testing"
)

func Test(t *testing.T) {
   filters, err := provider_filters()
   if err != nil {
      t.Fatal(err)
   }
   resp, err := filters[0].find()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("find.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
