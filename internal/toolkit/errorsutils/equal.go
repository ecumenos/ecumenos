package errorsutils

import "errors"

// Equals check if input values are equals
func Equals(left, right error) bool {
	if left == nil || right == nil {
		return errors.Is(left, right)
	}

	return left.Error() == right.Error()
}
