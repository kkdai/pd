package pd_test

import (
	"fmt"
	"time"

	. "github.com/kkdai/pubsub"
)

func ExamplePubsub() {
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
	fmt.Println(<-c2)
	//Got "test3"

	// Remove subscription "topic2" in c1
	ser.RemoveSubscription(c1, "topic2")

	// Publish new content in topic2
	ser.Publish("test4", "topic2")

	select {
	case val := <-c1:
		fmt.Errorf("Should not get %v notify on remove topic", val)
		break
	case <-time.After(time.Second):
		//Will go here, because we remove subscription topic2 in c1.
		fmt.Println("Not receive any msg from topic2, timeout.")
		break
	}
}

func ExamplePubsub_Subscribe() {
	ser := NewPubsub(1)
	//Subscribe topic1 in c1
	c1 := ser.Subscribe("topic1")
	fmt.Printf("c1 is a channel %v\n", c1)
	//Subscribe multiple topic "topic1", "topic2" and "topic3" in c2
	c2 := ser.Subscribe("topic1", "topic2", "topic3")
	fmt.Printf("c2 is a channel %v\n", c2)
}

func ExamplePubsub_AddSubscription() {
	ser := NewPubsub(1)
	//Subscribe topic1 in c1
	c1 := ser.Subscribe("topic1")
	//Add more Subscription "topic2" and "topic3" for c1
	ser.AddSubscription(c1, "topic2", "topic3")
}
