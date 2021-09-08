package fieldmask

import (
	"strings"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type FieldMask struct {
	paths []string
}

func (mask *FieldMask) Contains(field string) bool {
	if mask == nil {
		return true
	}

	for _, mask := range mask.paths {
		if strings.HasPrefix(field, mask) {
			return true
		}
	}

	return false
}

func NewFieldMask(mask *fieldmaskpb.FieldMask) *FieldMask {
	if len(mask.GetPaths()) == 0 {
		return nil
	}

	return &FieldMask{paths: mask.GetPaths()}
}
