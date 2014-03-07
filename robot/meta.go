package robot

import (
	"time"
)

type Meta struct {
	URL   string    `json:"url"`
	Added time.Time `json:"added"`
}
