package openRouter

import (
   "fmt"
   "log"
   "slices"
   "strings"
   "testing"
   "time"
)

func Test(t *testing.T) {
   filters, err := provider_filters()
   if err != nil {
      t.Fatal(err)
   }
   log.Println("len(filters)", len(filters))
   for i, filter := range filters {
      if i >= 1 {
         time.Sleep(99 * time.Millisecond)
      }
      models, err := filter.find()
      if err != nil {
         t.Fatal(err)
      }
      models = slices.DeleteFunc(models, func(m *model_data) bool {
         return strings.HasPrefix(m.Name, filter.Name + ": ")
      })
      fmt.Printf("%9v %v\n", len(models), filter.Name)
   }
}
