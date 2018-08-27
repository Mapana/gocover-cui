/*
	High satisfaction
	coverage: 100.0% of statements
*/
package example

import "testing"

func TestSUNHs(t *testing.T) {
	t.Run("group", func(t *testing.T) {
		t.Run("Test 1+1", func(t *testing.T) {
			if sun, err := SUN(1, 1); err != nil || sun != 2 {
				t.Errorf("SUN(%v, %v) = %v", 1, 1, sun)
			}
		})

		t.Run("Test 100+1000", func(t *testing.T) {
			if sun, err := SUN(100, 1000); err != nil || sun != 1100 {
				t.Errorf("SUN(%v, %v) = %v", 100, 1000, sun)
			}
		})

		t.Run("Test valid No parameters", func(t *testing.T) {
			if sun, err := SUN(); err.Error() != "at least two numbers" {
				t.Errorf("SUN() = %v", sun)
			}
		})

		t.Run("Test valid One parameters", func(t *testing.T) {
			if sun, err := SUN(1); err.Error() != "at least two numbers" {
				t.Errorf("SUN(%v) = %v", 1, sun)
			}
		})
	})
}