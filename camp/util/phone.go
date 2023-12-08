package utils

import "regexp"

// VerifyPhoneFormat 校验手机号
func VerifyPhoneFormat(phone string) bool {
	regular := `[0-9]+`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone) && len(phone) == 11
}

// PhoneReplace4 手机号屏蔽，前三后四
func PhoneReplace4(phone string) string {
	if phone == "" {
		return ""
	}
	l := len(phone)
	ph := "****"
	if l < 3 {
		return ph
	}
	if l < 7 {
		return phone[:3] + ph
	}
	return phone[:3] + ph + phone[7:]
}

// PhoneReplace6 手机号屏蔽，前三后二
func PhoneReplace6(phone string) string {
	if phone == "" {
		return ""
	}
	l := len(phone)
	ph := "******"
	if l < 3 {
		return ph
	}
	if l < 7 {
		return phone[:3] + ph
	}
	return phone[:3] + ph + phone[9:]
}
