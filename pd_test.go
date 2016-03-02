package pd

import (
	"log"
	"testing"
	"time"
)

func TestBasicFunction(t *testing.T) {
	ser := NewPubsub()
	c1 := ser.Subscribe("ch1")
	ser.Publish("test1", "ch1")

	if _, ok := <-c1; !ok {
		t.Error(" Error found on subscribed.\n")
	}
}

func TestTwoSubscribetor(t *testing.T) {
	ser := NewPubsub()
	c1 := ser.Subscribe("ch1")
	c2 := ser.Subscribe("ch2")

	ser.Publish("test2", "ch1")
	ser.Publish("test1", "ch2")

	val, ok := <-c1
	if !ok || string(val) != "test2" {
		t.Errorf("Error found \n")
	}

	val, ok = <-c2
	if !ok || string(val) != "test1" {
		t.Errorf("Error found \n")
	}
}

func TestAddSub(t *testing.T) {
	ser := NewPubsub()
	c1 := ser.Subscribe("ch1")
	ser.AddSubscription(c1, "ch2")
	ser.Publish("test2", "ch2")

	if val, ok := <-c1; !ok {
		t.Errorf("error on c1:%v", val)
	}
}

func TestRemoveSub(t *testing.T) {
	ser := NewPubsub()
	c1 := ser.Subscribe("ch1", "ch2")
	ser.Publish("test1", "ch2")

	if val, ok := <-c1; !ok {
		t.Errorf("error on addsub c1:%v", val)
	}

	ser.RemoveSubscription(c1, "ch1")
	ser.Publish("test2", "ch1")

	select {
	case val := <-c1:
		t.Errorf("Should not get %v notify on remove topic\n", val)
		break
	case <-time.After(time.Second):
		break
	}
}

func TestTopicList(t *testing.T) {
	ser := NewPubsub()
	t0 := ser.ListTopics()
	if len(t0) != 0 {
		t.Error("List error on empty")
	}

	c1 := ser.Subscribe("ch1", "ch2")
	t1 := ser.ListTopics()

	if len(t1) != 2 {
		t.Error("List error on basic count")
	}

	ser.AddSubscription(c1, "ch1", "ch2", "ch3")
	t2 := ser.ListTopics()
	if len(t2) != 3 {
		t.Error("List add error")
	}

	ser.RemoveSubscription(c1, "ch1")
	t3 := ser.ListTopics()
	if len(t3) != 2 {
		t.Error("List remove error")
	}
}

func BenchmarkAddSub(b *testing.B) {
	big := NewPubsub()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Subscribe("1234567890")
	}
}

func BenchmarkRemoveSub(b *testing.B) {
	big := NewPubsub()
	var subChans []chan []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c1 := big.Subscribe("1234567890")
		subChans = append(subChans, c1)
	}

	b.ResetTimer()
	for _, v := range subChans {
		big.RemoveSubscription(v, "1234567890")
	}
}

func BenchmarkBasicFunction(b *testing.B) {
	ser := NewPubsub()
	c1 := ser.Subscribe("ch1")

	for i := 0; i < b.N; i++ {
		ser.Publish("test1", "ch1")

		if _, ok := <-c1; !ok {
			log.Println(" Error found on subscribed.")
		}
	}
}
