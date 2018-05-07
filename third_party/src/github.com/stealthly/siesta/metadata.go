package siesta

import (
	"fmt"
	"sync"
	"time"
)

type Metadata struct {
	connector      Connector
	metadataExpire time.Duration
	cache          map[string]*metadataEntry
	refreshLock    sync.Mutex
}

func NetMetadata(connector Connector, metadataExpire time.Duration) *Metadata {
	return &Metadata{
		connector:      connector,
		metadataExpire: metadataExpire,
		cache:          make(map[string]*metadataEntry),
	}
}

func (tmc *Metadata) Get(topic string) ([]int32, error) {
	cache := tmc.cache[topic]
	if cache == nil {
		err := tmc.Refresh([]string{topic})
		if err != nil {
			return nil, err
		}
	}

	cache = tmc.cache[topic]
	if cache != nil {
		if cache.timestamp.Add(tmc.metadataExpire).Before(time.Now()) {
			err := tmc.Refresh([]string{topic})
			if err != nil {
				return nil, err
			}
		}

		cache = tmc.cache[topic]
		if cache != nil {
			return cache.partitions, nil
		}
	}

	return nil, fmt.Errorf("Could not get topic metadata for topic %s", topic)
}

func (tmc *Metadata) Refresh(topics []string) error {
	tmc.refreshLock.Lock()
	defer tmc.refreshLock.Unlock()

	topicMetadataResponse, err := tmc.connector.GetTopicMetadata(topics)
	if err != nil {
		return err
	}

	for _, topicMetadata := range topicMetadataResponse.TopicsMetadata {
		partitions := make([]int32, 0)
		for _, partitionMetadata := range topicMetadata.PartitionsMetadata {
			partitions = append(partitions, partitionMetadata.PartitionID)
		}
		tmc.cache[topicMetadata.Topic] = newMetadataEntry(partitions)
	}

	return nil
}

type metadataEntry struct {
	partitions []int32
	timestamp  time.Time
}

func newMetadataEntry(partitions []int32) *metadataEntry {
	return &metadataEntry{
		partitions: partitions,
		timestamp:  time.Now(),
	}
}
