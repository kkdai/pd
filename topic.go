package pd

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	. "github.com/kkdai/diskqueue"
)

type Topic struct {
	Name      string
	dataQueue WorkQueue
	exitChan  chan int

	chanList []chan string
}

func NewTopic(tname string) *Topic {
	t := new(Topic)
	t.Name = tname
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s-%d", t.Name, time.Now().UnixNano()))
	if err != nil {
		log.Println("Cannot init disk queue:", err)
		return nil
	}
	t.dataQueue = NewDiskqueue(t.Name, tmpDir)
	t.exitChan = make(chan int)
	go t.inLoop()
	return t
}

func (t *Topic) SendDataToChans(data string) {
	t.dataQueue.Put([]byte(data))
}

func (t *Topic) AddChan(newChan chan string) {
	for _, v := range t.chanList {
		if newChan == v {
			//chan exist return
			return
		}
	}
	t.chanList = append(t.chanList, newChan)
}

func (t *Topic) Cleanup() {
	//cleanup
	t.exitChan <- 1
}

func (t *Topic) RemoveChan(delChan chan string) {
	t.chanList = removeChanFromSlice(t.chanList, delChan)
}

func (t *Topic) CountChanList() int {
	return len(t.chanList)
}

func removeChanFromSlice(slice []chan string, target chan string) []chan string {
	var retSlice []chan string
	removeIndex := -1
	for k, v := range slice {
		if v == target {
			removeIndex = k
		}
	}

	if removeIndex == -1 {
		return slice
	}
	if len(slice) == 1 && removeIndex == 0 {
		return retSlice
	} else {
		retSlice = append(slice[:removeIndex], slice[removeIndex+1:]...)
	}
	return retSlice
}

func (t *Topic) inLoop() {
	var dataRead []byte
	for {
		dataRead = nil
		select {
		case dataRead = <-t.dataQueue.ReadChan():
			log.Println("TOPIC:Got data ", string(dataRead))
			//Got data start to put on all chan
			for _, targetChan := range t.chanList {
				targetChan <- string(dataRead)
			}
		case <-t.exitChan:
			//exist
			goto exit
		}

	}

exit:
	//cleanup
	t.dataQueue.Empty()
	t.dataQueue.Close()
}
