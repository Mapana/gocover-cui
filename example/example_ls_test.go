/*
	Low satisfaction
	coverage: 85.7% of statements
*/
package example

import "testing"

func TestSUNLs(t *testing.T) {
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
	})
}
