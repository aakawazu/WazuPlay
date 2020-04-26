package db

import (
	"fmt"
	"time"
)

// TimeNow return timestamp
func TimeNow(w time.Duration) string {
	t := fmt.Sprintf("%s",
		(time.Now().Add(w * time.Minute)).Round(time.Second),
	)
	t = t[0 : len(t)-4]
	return t
}
