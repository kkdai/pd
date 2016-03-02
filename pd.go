package pd

import "sync"

// PD :Pubsub struct Only content a userIndex and accessDB which content a chan map
type PD struct {
	sync.RWMutex

	//To quick refer if topic exist
	topicMap map[string]int

	//map to store "topic -> chan List" for publish
	topicMapClients map[string]*Topic
}

//Subscribe : Subscribe channels, the channels could be a list of channels name
//     The channel name could be any, without define in server
func (p *PD) Subscribe(topics ...string) chan []byte {
	p.RLock()
	defer p.RUnlock()

	//init new chan using capacity as channel buffer
	workChan := make(chan []byte)
	p.updateTopicMapClient(workChan, topics)
	return workChan
}

//ListTopics :Return all the topic
func (p *PD) ListTopics() []string {
	p.RLock()
	defer p.RUnlock()

	var retSlice []string
	for k := range p.topicMap {
		retSlice = append(retSlice, k)
	}
	return retSlice
}

func (p *PD) updateTopicMapClient(clientChan chan []byte, topics []string) {
	for _, topic := range topics {
		if _, exist := p.topicMap[topic]; !exist {
			p.topicMap[topic] = 1
			p.topicMapClients[topic] = NewTopic(topic)
		}
		p.topicMapClients[topic].AddChan(clientChan)
	}
}

//AddSubscription  Add a new topic subscribe to specific client channel.
func (p *PD) AddSubscription(clientChan chan []byte, topics ...string) {
	p.RLock()
	defer p.RUnlock()

	p.updateTopicMapClient(clientChan, topics)
}

//RemoveSubscription Remove sub topic list on specific chan
func (p *PD) RemoveSubscription(clientChan chan []byte, topics ...string) {
	p.RLock()
	defer p.RUnlock()

	for _, topic := range topics {
		//Remove from topic->chan map
		if topicObj, ok := p.topicMapClients[topic]; ok {
			//remove one client chan in chan List
			topicObj.RemoveChan(clientChan)

			if topicObj.CountChanList() == 0 {
				//Don't have any subscription
				topicObj.Cleanup()
				delete(p.topicMap, topic)
				delete(p.topicMapClients, topic)
			}
		}

	}
}

//Publish  Publish a content to a list of channels
//         The content could be any type.
func (p *PD) Publish(content string, topics ...string) {
	p.RLock()
	defer p.RUnlock()

	for _, topic := range topics {
		if topicObj, ok := p.topicMapClients[topic]; ok {
			topicObj.SendDataToChans(content)
		}
	}
}

// NewPubsub :Create a pubsub with expect init size, but the size could be extend.
func NewPubsub() *PD {
	//log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)
	server := PD{
		topicMap:        make(map[string]int),
		topicMapClients: make(map[string]*Topic),
	}
	return &server
}
