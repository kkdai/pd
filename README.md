PD: Pubsub service base on disk queue on golang
==============

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/pd/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/pd?status.svg)](https://godoc.org/github.com/kkdai/pd)  [![Build Status](https://travis-ci.org/kkdai/PD.svg?branch=master)](https://travis-ci.org/kkdai/pd)



What is Pubsub
=============
Pubsub is prove of concept implement for [Redis](http://redis.io/) "Pub/Sub" messaging management feature. SUBSCRIBE, UNSUBSCRIBE and PUBLISH implement the Publish/Subscribe messaging paradigm where (citing Wikipedia) senders (publishers) are not programmed to send their messages to specific receivers (subscribers).  (sited from [here](http://redis.io/topics/pubsub))

What is Disk Queue
=============
[Disk Queue](https://github.com/nsqio/nsq/blob/master/nsqd/diskqueue.go) is a data structure which come from [NSQ](https://github.com/nsqio/nsq). It is message queue data structure using disk file as storage medium.

How it work together
=============

This is no API change from Pub/Sub mechanism, but change it basic concurrency process.

`Topic` as a another object to handle Topic related info with Disk Queue. Each `Topic` contains two kind information:

- `Data Queue`: Which is all publish data send to this topic.
- `Channel List`: Which is who subscribe this topic, we need notify.

In this modification, it gain follow benefits:

- Infinite Topic Queue Size (depends on storage size)
- Buffering Publish to improve performance.



Installation and Usage
=============


Install
---------------
```
go get github.com/kkdai/pd
```

Usage
---------------

```go
package main
    
import (
	"fmt"
    
	. "github.com/kkdai/pubsub"
)
    
func main() {
	ser := NewPubsub(1)
	c1 := ser.Subscribe("topic1")
	c2 := ser.Subscribe("topic2")
	ser.Publish("test1", "topic1")
	ser.Publish("test2", "topic2")
	fmt.Println(<-c1)
	//Got "test1"
	fmt.Println(<-c2)
	//Got "test2"


    // Add subscription "topic2" for c1.          
	ser.AddSubscription(c1, "topic2")

    // Publish new content in topic2
	ser.Publish("test3", "topic2")

	fmt.Println(<-c1)
	//Got "test3"
	
    // Remove subscription "topic2" in c1
	ser.RemoveSubscription(c1, "topic2")
	
    // Publish new content in topic2
	ser.Publish("test4", "topic2")
    
	select {
	case val := <-c1:
		fmt.Printf("Should not get %v notify on remove topic\n", val)
		break
	case <-time.After(time.Second):
	    //Will go here, because we remove subscription topic2 in c1.         
		break
	}
}
```

Benchmark
---------------

Benchmark include memory usage. (original memory)

```
BenchmarkAddSub-4       	     500	2906467 ns/op
BenchmarkRemoveSub-4    	   10000	 232910 ns/op
BenchmarkBasicFunction-4	 5000000	    232 ns/op
```

Benchmark include memory usage. (Using Disk Queue)

```
BenchmarkAddSub-4       	  300000	125628 ns/op 
BenchmarkRemoveSub-4    	  200000    144854 ns/op
BenchmarkBasicFunction-4	    2000	906076 ns/op
```

Inspired By
---------------

- [Redis: Pubsub](http://redis.io/topics/pubsub)
- [NSQ: Disqueue](https://github.com/nsqio/nsq/blob/master/nsqd/diskqueue.go)
- [NSQ: realtime distributed message processing at scale](http://word.bitly.com/post/33232969144/nsq)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.
