// Package retry retries an operation with exponential backoff.
package retry

import "time"

// Do runs fn up to attempts times. After each failure it waits, doubling the
// delay each round (exponential backoff), and returns the final error if every
// attempt fails.
func Do(attempts int, base time.Duration, fn func() error) error {
      delay := base
      var err error
      for i := 0; i < attempts; i++ {
              if err = fn(); err == nil {
                      return nil
              }
              if i < attempts-1 {
                      time.Sleep(delay)
                      delay *= 2
              }
      }
      return err
}
