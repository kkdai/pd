package pd

import (
	. "github.com/kkdai/diskqueue"
)

type chanMapStringList map[chan interface{}][]string
type stringMapChanList map[string][]chan interface{}

// Pubsub struct: Only content a userIndex and accessDB which content a chan map
type PD struct {
	//Capacity for each chan buffer
	capacity int

	//map to store "chan -> Topic List" for find subscription
	clientMapTopics chanMapStringList
	//map to store "topic -> chan List" for publish
	topicMapClients stringMapChanList
}

//Sub: Subscribe channels, the channels could be a list of channels name
//     The channel name could be any, without define in server
func (p *PD) Subscribe(topics ...string) chan interface{} {
	//init new chan using capacity as channel buffer
	workChan := make(chan interface{}, p.capacity)
	p.updateTopicMapClient(workChan, topics)
	return workChan
}

func (p *PD) updateTopicMapClient(clientChan chan interface{}, topics []string) {
	var updateChanList []chan interface{}
	for _, topic := range topics {
		updateChanList, _ = p.topicMapClients[topic]
		updateChanList = append(updateChanList, clientChan)
		p.topicMapClients[topic] = updateChanList
	}
	p.clientMapTopics[clientChan] = topics
}

//AddSubscription:  Add a new topic subscribe to specific client channel.
func (p *PD) AddSubscription(clientChan chan interface{}, topics ...string) {
	p.updateTopicMapClient(clientChan, topics)
}

//RemoveSubscription: Remove sub topic list on specific chan
func (p *PD) RemoveSubscription(clientChan chan interface{}, topics ...string) {

	for _, topic := range topics {
		//Remove from topic->chan map
		if chanList, ok := p.topicMapClients[topic]; ok {
			//remove one client chan in chan List
			var updateChanList []chan interface{}
			for _, client := range chanList {
				if client != clientChan {
					updateChanList = append(updateChanList, client)
				}
			}
			p.topicMapClients[topic] = updateChanList
		}

		//Remove from chan->topic map
		if topicList, ok := p.clientMapTopics[clientChan]; ok {
			var updateTopicList []string
			for _, updateTopic := range topicList {
				if updateTopic != topic {
					updateTopicList = append(updateTopicList, topic)
				}
			}
			p.clientMapTopics[clientChan] = updateTopicList
		}
	}
}

//Publish: Publish a content to a list of channels
//         The content could be any type.
func (p *PD) Publish(content interface{}, topics ...string) {
	for _, topic := range topics {
		if chanList, ok := p.topicMapClients[topic]; ok {
			//Someone has subscribed this topic
			for _, channel := range chanList {
				channel <- content
			}
		}
	}
}

// Create a pubsub with expect init size, but the size could be extend.
func NewPubsub(initChanCapacity int) *PD {
	initClientMapTopics := make(chanMapStringList)
	initTopicMapClients := make(stringMapChanList)

	server := PD{clientMapTopics: initClientMapTopics, topicMapClients: initTopicMapClients}
	server.capacity = initChanCapacity
	return &server
}
