package fromproto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Timestamp(in *timestamppb.Timestamp) time.Time {
	if in == nil {
		return time.Time{}
	}

	return in.AsTime()
}
