// Package parser provides functionality for parsing standard 5-field cron
// expressions into structured representations.
//
// A cron expression consists of five space-separated fields:
//
//	┌─────────── minute        (0–59)
//	│ ┌───────── hour          (0–23)
//	│ │ ┌─────── day of month  (1–31)
//	│ │ │ ┌───── month         (1–12)
//	│ │ │ │ ┌─── day of week   (0–6, Sunday=0)
//	│ │ │ │ │
//	* * * * *
//
// Supported syntax:
//   - *        — wildcard (every value in range)
//   - N        — exact value
//   - N-M      — inclusive range
//   - */S      — every S steps across the full range
//   - N-M/S    — every S steps within range N-M
//   - A,B,...  — comma-separated list of any of the above
//
// Example usage:
//
//	expr, err := parser.Parse("0 9-17 * * 1-5")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(expr.Hour.Values) // [9 10 11 12 13 14 15 16 17]
package parser
