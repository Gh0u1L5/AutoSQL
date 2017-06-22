package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SQLmapTask struct {
	id     string
	server string
}

type SQLmapTaskInfo struct {
	Url string `json:"url"`
	// Headers map[string]interface{}
}

func (t SQLmapTask) unmarshalJSON(resp *http.Response, ptr interface{}) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if err := json.Unmarshal(data, ptr); err != nil {
		log.Panic(err)
	}
}

func (t *SQLmapTask) create(server string) {
	if server != "" {
		t.server = server
	} else {
		t.server = "http://127.0.0.1:8775"
	}

	resp, err := http.Get(t.server + "/task/new")
	if err != nil {
		log.Panic(err)
	}
	var result struct {
		Success bool
		Taskid  string
	}
	t.unmarshalJSON(resp, &result)
	if !result.Success {
		log.Panic("Failed to create new task")
	}
	t.id = result.Taskid
}

func (t *SQLmapTask) delete() {
	resp, err := http.Get(t.server + "/task/" + t.id + "/delete")
	if err != nil {
		log.Panic(err)
	}
	var result struct {
		Success bool
		Message string
	}
	t.unmarshalJSON(resp, &result)
	if !result.Success {
		log.Panic(result.Message)
	}
}

func (t *SQLmapTask) start(info SQLmapTaskInfo) {
	data, err := json.Marshal(info)
	if err != nil {
		log.Panic(err)
	}
	resp, err := http.Post(t.server+"/scan/"+t.id+"/start", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
	var result struct {
		Success  bool
		EngineID int
	}
	t.unmarshalJSON(resp, &result)
	if !result.Success {
		log.Panic("Failed to start task " + t.id)
	}
}

func (t SQLmapTask) status() string {
	resp, err := http.Get(t.server + "/scan/" + t.id + "/status")
	if err != nil {
		log.Panic(err)
	}
	var result struct {
		Success bool
		Status  string
	}
	t.unmarshalJSON(resp, &result)
	if !result.Success {
		log.Panic("Failed to check status of task " + t.id)
	}
	return result.Status
}

func (t SQLmapTask) result() interface{} {
	resp, err := http.Get(t.server + "/scan/" + t.id + "/data")
	if err != nil {
		log.Panic(err)
	}
	var result struct {
		Success bool
		Data    interface{}
	}
	t.unmarshalJSON(resp, &result)
	if !result.Success {
		log.Panic("Failed to get result of task " + t.id)
	}
	return result.Data
}
