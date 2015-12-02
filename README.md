# ensqs
The package to enqueue to AWS SQS.

###Install
```
go get github.com/LeeQY/ensqs
```

###Introduction
* If error happens when enqueue to SQS, users need to retry.
* If the network is poor, it will be long to return.
* If there are many jobs needed to be enqueued, batch methods need to be applied.

This is the package to handle the problems above.

Enqueue to SQS is handled in a seperate coroutine. And will be cached in memory, which result in fast return. Automatically use batch method.

