package veritrans

import "time"

func getAfterOneMonth() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(0, 1, 0)
	return expiredAt.Format("01/06")
}

func getAfterOneYear() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(1, 0, 0)
	return expiredAt.Format("01/06")
}
