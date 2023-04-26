package postbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpec_Get(t *testing.T) {
	tests := []struct {
		name string
		spec *SubjectSpec
		want []string
	}{
		{"foo", NewSubjectSpec().AddPath("foo"), []string{"foo"}},
		{"foo.bar", NewSubjectSpec().AddMulti("foo.bar"), []string{"foo", "foo.bar"}},
		{"foo.bar single", NewSubjectSpec().AddPath("foo.bar"), []string{"foo.bar"}},
		{"foo.bar.baz exclude", NewSubjectSpec().AddMulti("foo.bar.baz").Exclude("foo"), []string{"foo.bar", "foo.bar.baz"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.Get()
			assert.ElementsMatch(t, tt.want, result)
		})
	}
}
