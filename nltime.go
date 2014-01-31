package nltime

import (
   "errors"
   "strings"
   "time"
)

func ParseRange(s string, firstDayOfWeek time.Weekday) ([]time.Time, error) {
   mod:= 0
   start := time.Now()
   end := start
   monthSet := false

   toks := strings.Split(s, " ")
   tokcnt := len(toks)
   for _, tok := range toks {
      switch tok {
      case "last":
         if tokcnt < 2 {
            return []time.Time{}, errors.New("last what? year? month? week?")
         }
         mod = -1
      case "this":
         if tokcnt < 2 {
            return []time.Time{}, errors.New("this what? year? month? week?")
         }
         mod = 0
      case "next":
         if tokcnt < 2 {
            return []time.Time{}, errors.New("next what? year? month? week?")
         }
         mod = 1
      case "year":
         loc := time.Now().Location()
         if monthSet {
            start = time.Date(start.Year() + mod, start.Month(), 1, 0, 0, 0, 0, loc)
            end = time.Date(end.Year() + mod, start.Month(), LastDayOfMonth(start), 0, 0, 0, 0, loc)
         } else {
            start = time.Date(start.Year() + mod, time.January, 1, 0, 0, 0, 0, loc)
            end = time.Date(end.Year() + mod, time.December, 31, 0, 0, 0, 0, loc)
         }
      case "month":
         start = start.AddDate(0, mod, -(start.Day() - 1))
         end = end.AddDate(0, mod, 0)
         end = end.AddDate(0, 0, LastDayOfMonth(end) - end.Day())
         monthSet = true
      case "week":
         days := 7 * mod
         fdow := int(firstDayOfWeek)
         wkday := int(start.Weekday())
         start = start.AddDate(0, 0, days - (wkday - fdow))
         end = end.AddDate(0, 0, days + (7 - wkday))
      }
   }
   return []time.Time{start,end}, nil
}

func LastDayOfMonth(t time.Time) int {
   switch t.Month() {
   case time.February:
      if IsLeap(t.Year()) {
         return 29
      }
      return 28
   case time.April, time.June, time.September, time.November:
      return 30
   }
   return 31
}

func IsLeap(year int) bool {
   return year%4 == 0 && (year%10 != 0 || year%400 == 0)
}
