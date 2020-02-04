//MIT License
//
//Copyright (c) 2020 Wout Slakhorst
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MopidyClient struct {
	RPCAddress string
}

// clearPlaylist clears the current playlist
//{
//  "jsonrpc": "2.0",
//  "id": 1,
//  "method": "core.tracklist.clear"
//}
func (mc MopidyClient) clearPlaylist() error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.tracklist.clear",
	})

	return err
}

// shufflePlaylist shuffles the current playlist
//{
//  "jsonrpc": "2.0",
//  "id": 1,
//  "method": "core.tracklist.shiffle"
//}
func (mc MopidyClient) shufflePlaylist() error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.tracklist.shuffle",
	})

	return err
}

// setTracklist adds the given album/playlist to the tracklist
//{
//  "jsonrpc": "2.0",
//  "id": 1,
//  "method": "core.tracklist.add",
//  "params": {
//    "uris": ["spotify:playlist:5vuSj9MXB11ZiLulJjx0Ag"]
//  }
//}
func (mc MopidyClient) setTracklist(url string) error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.tracklist.add",
		Params: map[string]interface{}{
			"uris": []string{
				url,
			},
		},
	})

	return err
}

// play starts playback
//{
//"jsonrpc": "2.0",
//"id": 1,
//"method": "core.playback.play"
//}
func (mc MopidyClient) play() error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.playback.play",
	})

	return err
}

// stop stops playback
//{
//"jsonrpc": "2.0",
//"id": 1,
//"method": "core.playback.play"
//}
func (mc MopidyClient) stop() error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.playback.stop",
	})

	return err
}

// next plays the next track
//{
//"jsonrpc": "2.0",
//"id": 1,
//"method": "core.playback.next"
//}
func (mc MopidyClient) next() error {
	_, err := mc.do(MopidyRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "core.playback.next",
	})

	return err
}

type MopidyRequest struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Id      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type MopidyResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	Error   MopidyError `json:"error"`
}

type MopidyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (mc MopidyClient) do(req MopidyRequest) (*http.Response, error) {
	fmt.Printf("%v\n", req)
	j, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := bytes.NewReader(j)

	fmt.Println(string(j))

	resp, err := http.Post(mc.RPCAddress, "application/json", reader)

	b, _ := ioutil.ReadAll(resp.Body)
	mopidyResponse := MopidyResponse{}
	json.Unmarshal(b, &mopidyResponse)

	if mopidyResponse.Error.Code != 0 {
		fmt.Println(string(b))
	}

	return resp, err
	//return http.Post(mc.RPCAddress, "application/json", reader)
}
