package ensqs

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestEnsqs(t *testing.T) {

	thisCount := 1000
	for i := 0; i < thisCount; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		num := r.Intn(200)

		s := strconv.Itoa(i)
		v := Value{&s, []byte{}}

		time.Sleep(time.Duration(num) * time.Microsecond)

		AddValue(&v)
	}
	time.Sleep(1 * time.Second)

	if int32(thisCount) != count {
		t.Error("Error in handle values. count: ", count)
	}
}
