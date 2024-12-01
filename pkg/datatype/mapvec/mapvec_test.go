package mapvec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestICanAddInMapVec(t *testing.T) {
	trial1 := New[string, string]().
		Add("salut", "les kids").
		Add("donde", "   estaaaaaaa ").
		Add("donde", "en la ").
		Add("donde", "pl4y4$1 ")
	goal1 := MapVec[string, string]{
		"salut": []string{"les kids"},
		"donde": []string{"   estaaaaaaa ", "en la ", "pl4y4$1 "},
	}
	assert.Equal(t, goal1, trial1)
	type dummy struct {
		t float32
	}
	trial2 := New[float32, dummy]().
		Add(1.12, dummy{1.12}).
		Add(3.231321312312, dummy{3.231321312312})
	goal2 := MapVec[float32, dummy]{
		1.12:           []dummy{{1.12}},
		3.231321312312: []dummy{{3.231321312312}},
	}
	assert.Equal(t, goal2, trial2)
}
