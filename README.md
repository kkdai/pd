PD: Pubsub service base on disk queue on golang
==============

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/PD/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/PD?status.svg)](https://godoc.org/github.com/kkdai/PD)  [![Build Status](https://travis-ci.org/kkdai/PD.svg?branch=master)](https://travis-ci.org/kkdai/PD)



What is Pubsub
=============
Pubsub is prove of concept implement for [Redis](http://redis.io/) "Pub/Sub" messaging management feature. SUBSCRIBE, UNSUBSCRIBE and PUBLISH implement the Publish/Subscribe messaging paradigm where (citing Wikipedia) senders (publishers) are not programmed to send their messages to specific receivers (subscribers).  (sited from [here](http://redis.io/topics/pubsub))


Installation and Usage
=============


Install
---------------
        go get github.com/kkdai/pubsub


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

Benchmark include memory usage.

```
BenchmarkAddSub-4       	     500	   2906467 ns/op	 1605949 B/op	       3 allocs/op
BenchmarkRemoveSub-4    	   10000	    232910 ns/op	  174260 B/op	      16 allocs/op
BenchmarkBasicFunction-4	 5000000	       232 ns/op	      19 B/op	       1 allocs/op
```

Inspired By
---------------


- [Redis: Pubsub](http://redis.io/topics/pubsub)
- [chandru/pubsub](https://github.com/tuxychandru/pubsub)


Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.
