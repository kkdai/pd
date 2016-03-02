package pd

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	diskqueue "github.com/kkdai/diskqueue"
)

//Topic : Topic strucuture
type Topic struct {
	Name      string
	dataQueue diskqueue.WorkQueue
	exitChan  chan int

	chanList []chan []byte
}

//NewTopic :Create new topic with a disk queue
func NewTopic(tname string) *Topic {
	t := new(Topic)
	t.Name = tname
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s-%d", t.Name, time.Now().UnixNano()))
	if err != nil {
		log.Println("Cannot init disk queue:", err)
		return nil
	}
	t.dataQueue = diskqueue.NewDiskqueue(t.Name, tmpDir)
	t.exitChan = make(chan int)
	go t.inLoop()
	return t
}

//SendDataToChans :Send data to the subscriptor
func (t *Topic) SendDataToChans(data []byte) {
	t.dataQueue.Put(data)
}

//AddChan :Add new subscriptor this is topic
func (t *Topic) AddChan(newChan chan []byte) {
	for _, v := range t.chanList {
		if newChan == v {
			//chan exist return
			return
		}
	}
	t.chanList = append(t.chanList, newChan)
}

//Cleanup :Close and exist inLoop()
func (t *Topic) Cleanup() {
	//cleanup
	t.exitChan <- 1
}

//RemoveChan :Remove a subscriptor from this topic
func (t *Topic) RemoveChan(delChan chan []byte) {
	t.chanList = removeChanFromSlice(t.chanList, delChan)
}

//CountChanList :Count the all subscriptors
func (t *Topic) CountChanList() int {
	return len(t.chanList)
}

func removeChanFromSlice(slice []chan []byte, target chan []byte) []chan []byte {
	var retSlice []chan []byte
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
	}
	retSlice = append(slice[:removeIndex], slice[removeIndex+1:]...)
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
				targetChan <- dataRead
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
