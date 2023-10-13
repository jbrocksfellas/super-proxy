package main

import "sync"

type TrafficStore struct {
	totalTraffic    int64
	userTraffic     map[string]int64
	trafficMu       sync.Mutex
	maxTotalTraffic int64
	maxUserTraffic  int64
}

func NewTrafficStore(maxTotalTraffic int64, maxUserTraffic int64) *TrafficStore {
	return &TrafficStore{
		totalTraffic:    0,
		trafficMu:       sync.Mutex{},
		userTraffic:     map[string]int64{},
		maxTotalTraffic: maxTotalTraffic * 1000000,
		maxUserTraffic:  maxUserTraffic * 1000000,
	}
}

func (t *TrafficStore) addTraffic(username string, traffic int64) {
	t.trafficMu.Lock()
	defer t.trafficMu.Unlock()

	// add user traffic
	_, ok := t.userTraffic[username]
	if !ok {
		t.userTraffic[username] = traffic
		t.totalTraffic += traffic
	} else {
		t.userTraffic[username] += traffic
		t.totalTraffic += traffic
	}

}

func (t *TrafficStore) getUserTraffic(usrname string) int64 {
	t.trafficMu.Lock()
	defer t.trafficMu.Unlock()

	val, ok := t.userTraffic[usrname]
	if !ok {
		return 0
	}

	return val
}

func (t *TrafficStore) getTotalTraffic() int64 {
	t.trafficMu.Lock()
	defer t.trafficMu.Unlock()

	return t.totalTraffic
}
