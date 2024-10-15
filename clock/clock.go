// time.Timeだと精密すぎてテストデータの比較に困る。
// 時刻を固定化するためのpkg

package clock

import "time"

type Clocker interface {
	Now() time.Time
}

// 以下の二つはClocker抽象型を満たす
// どの具体的な Clocker 型（RealClocker または FixedClocker）を渡すかによって決まる

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (f FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
