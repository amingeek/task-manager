// backend/controllers/controller_utils.go
package controllers

import (
	"strconv"
)

// توابع کمکی مشترک بین تمام کنترلرها

// StringToUint - تبدیل ایمن string به uint
func StringToUint(s string) uint {
	if s == "" {
		return 0
	}
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(num)
}

// StringToUintWithDefault - تبدیل با مقدار پیش‌فرض
func StringToUintWithDefault(s string, defaultValue uint) uint {
	if s == "" {
		return defaultValue
	}
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return uint(num)
}
