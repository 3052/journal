package main

import (
   "fmt"
   "strconv"
   "strings"
   "time"
)

var events = []jog{
   {
      days_before: 8,
      day_of: date(2026, 2, 6),
      ok: nil,
   },
   {
      days_before: 9,
      day_of: date(2026, 1, 29),
      ok: some(true),
   },
   {
      days_before: 10,
      day_of: date(2026, 1, 20),
      ok: some(true),
   },
   {
      days_before: 11,
      day_of: date(2026, 1, 10),
      ok: some(true),
   },
   {
      days_before: 10,
      day_of: date(2025, 12, 26),
      ok: some(false),
   },
   {
      days_before: 9,
      day_of: date(2025, 12, 16),
      ok: some(false),
   },
   {
      days_before: 8,
      day_of: date(2025, 12, 7),
      ok: some(false),
   },
   {
      days_before: 7,
      day_of: date(2025, 11, 29),
      ok: some(false),
   },
}

func main() {
   for i, event := range events {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&event)
   }
}

func date(year int, month time.Month, day int) time.Time {
   return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func some(value bool) *bool {
   return &value
}

func (j *jog) String() string {
   var data strings.Builder
   data.WriteString("days before = ")
   data.WriteString(strconv.Itoa(j.days_before))
   data.WriteString("\nday of = ")
   data.WriteString(j.day_of.Format(time.DateOnly))
   if j.ok != nil {
      data.WriteString("\nok = ")
      data.WriteString(strconv.FormatBool(*j.ok))
   } else {
      data.WriteString("\nok = nil")
   }
   return data.String()
}

type jog struct {
   days_before int
   day_of      time.Time
   ok          *bool
}
