package ensqs

type Value struct {
	key string
	v   []byte
}

var (
	qInfo      string
	regionInfo string
)

func SetInfo(q, region *string) {
	if q == nil || region == nil {
		panic("Info can't be nil.")
	}
	qInfo = *q
	regionInfo = *region
}
