package ensqs

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type InfoQ struct {
	m          sync.RWMutex
	values     []Value
	processing bool
}

type Value struct {
	Key *string
	V   []byte
}

var (
	qURL string

	infoQ InfoQ
	//	count int32

	thisSQS   *sqs.SQS
	waitTime  = int64(20)
	binary    = "Binary"
	valueName = "value"
)

// Need to set once.
func SetInfo(queueURL, region *string) {
	if queueURL == nil || region == nil {
		panic("Info can't be nil.")
	}
	qURL = *queueURL

	thisSQS = sqs.New(session.New(&aws.Config{Region: region}))
}

func AddValue(v *Value) error {
	if v == nil {
		return errors.New("Value can't be nil.")
	}
	infoQ.m.Lock()
	defer infoQ.m.Unlock()

	infoQ.values = append(infoQ.values, *v)
	if !infoQ.processing {
		infoQ.processing = true
		go send()
	}
	return nil
}

func getValues() []Value {
	infoQ.m.RLock()
	defer infoQ.m.RUnlock()

	if len(infoQ.values) > 10 {
		return infoQ.values[:10]
	} else {
		return infoQ.values
	}
}

func deleteValue(l int) {
	infoQ.m.Lock()
	defer infoQ.m.Unlock()

	infoQ.values = infoQ.values[l:]
}

func send() {
	values := getValues()
	l := len(values)

	if l == 0 {
		infoQ.m.Lock()
		infoQ.processing = false
		infoQ.m.Unlock()
		return
	}

	entries := make([]*sqs.SendMessageBatchRequestEntry, l)
	for i := 0; i < l; i++ {
		s := strconv.Itoa(i)

		var one sqs.SendMessageBatchRequestEntry
		one.Id = &s
		one.MessageBody = values[i].Key

		var mav sqs.MessageAttributeValue
		mav.BinaryValue = values[i].V
		mav.DataType = &binary

		one.MessageAttributes = make(map[string]*sqs.MessageAttributeValue)
		one.MessageAttributes[valueName] = &mav

		entries[i] = &one
	}

	params := &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: &qURL,
	}

	if _, err := thisSQS.SendMessageBatch(params); err != nil {
		log.Println(err)
	} else {
		deleteValue(l)
		//	atomic.AddInt32(&count, int32(l))
	}

	send()
}
