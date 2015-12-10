package ensqs

import (
	"log"
	"strconv"

	"github.com/LeeQY/jobber"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Value struct {
	Key *string
	V   []byte
}

var (
	qURL string

	thisSQS   *sqs.SQS
	waitTime  = int64(20)
	binary    = "Binary"
	valueName = "value"

	j *jobber.Jobber
)

// Need to set once.
func SetInfo(queueURL, region *string) {
	if queueURL == nil || region == nil {
		panic("Info can't be nil.")
	}
	qURL = *queueURL

	thisSQS = sqs.New(session.New(&aws.Config{Region: region}))

	j = jobber.New(batchSend, 10)
}

func AddValue(v *Value) {
	j.AddJob(v)
}

func batchSend(values []interface{}) bool {
	l := len(values)
	entries := make([]*sqs.SendMessageBatchRequestEntry, l)

	for i := 0; i < l; i++ {
		s := strconv.Itoa(i)

		value := values[i].(*Value)

		var one sqs.SendMessageBatchRequestEntry
		one.Id = &s
		one.MessageBody = (*value).Key

		var mav sqs.MessageAttributeValue
		mav.BinaryValue = (*value).V
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
		return false
	}
	return true
}
