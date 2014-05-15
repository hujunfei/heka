/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012-2014
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Rob Miller (rmiller@mozilla.com)
#   Mike Trinkala (trink@mozilla.com)
#
# ***** END LICENSE BLOCK *****/

package tcp

import (
	"code.google.com/p/gomock/gomock"
	. "github.com/mozilla-services/heka/pipeline"
	pipeline_ts "github.com/mozilla-services/heka/pipeline/testsupport"
	plugins_ts "github.com/mozilla-services/heka/plugins/testsupport"
	gs "github.com/rafrombrc/gospec/src/gospec"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func TcpOutputSpec(c gs.Context) {
	t := new(pipeline_ts.SimpleT)
	ctrl := gomock.NewController(t)

	tmpDir, tmpErr := ioutil.TempDir("", "tcp-tests-")
	os.MkdirAll(tmpDir, 0777)
	defer func() {
		ctrl.Finish()
		tmpErr = os.RemoveAll(tmpDir)
		c.Expect(tmpErr, gs.Equals, nil)
	}()
	pConfig := NewPipelineConfig(nil)

	c.Specify("TcpOutput", func() {
		origBaseDir := Globals().BaseDir
		Globals().BaseDir = tmpDir
		defer func() {
			Globals().BaseDir = origBaseDir
		}()

		tcpOutput := new(TcpOutput)
		tcpOutput.SetName("test")
		config := tcpOutput.ConfigStruct().(*TcpOutputConfig)
		tcpOutput.Init(config)

		tickChan := make(chan time.Time)
		oth := plugins_ts.NewOutputTestHelper(ctrl)
		oth.MockOutputRunner.EXPECT().Ticker().Return(tickChan)
		encoder := new(ProtobufEncoder)
		encoder.Init(nil)
		oth.MockOutputRunner.EXPECT().Encoder().Return(encoder)

		var wg sync.WaitGroup
		inChan := make(chan *PipelinePack, 1)

		msg := pipeline_ts.GetTestMessage()
		pack := NewPipelinePack(pConfig.InputRecycleChan())
		pack.Message = msg
		pack.Decoded = true

		outStr := "Write me out to the network"
		matchBytes := make([]byte, 0, 1000)
		newpack := NewPipelinePack(nil)
		newpack.Message = msg
		newpack.Decoded = true
		newpack.Message.SetPayload(outStr)
		err := ProtobufEncodeMessage(newpack, &matchBytes)
		c.Expect(err, gs.IsNil)

		c.Specify("writes out to the network", func() {
			inChanCall := oth.MockOutputRunner.EXPECT().InChan().AnyTimes()
			inChanCall.Return(inChan)

			collectData := func(ch chan string) {
				ln, err := net.Listen("tcp", "localhost:9125")
				if err != nil {
					ch <- err.Error()
					return
				}
				ch <- "ready"
				conn, err := ln.Accept()
				if err != nil {
					ch <- err.Error()
					return
				}
				b := make([]byte, 1000)
				n, _ := conn.Read(b)
				ch <- string(b[0:n])
				conn.Close()
				ln.Close()
			}
			ch := make(chan string, 1) // don't block on put
			go collectData(ch)
			result := <-ch // wait for server

			err := tcpOutput.Init(config)
			c.Assume(err, gs.IsNil)

			pack.Message.SetPayload(outStr)
			go func() {
				wg.Add(1)
				c.Expect(err, gs.IsNil)
				err = tcpOutput.Run(oth.MockOutputRunner, oth.MockHelper)
				wg.Done()
			}()

			msgcount := atomic.LoadInt64(&tcpOutput.processMessageCount)
			c.Expect(msgcount, gs.Equals, int64(0))

			inChan <- pack
			result = <-ch

			msgcount = atomic.LoadInt64(&tcpOutput.processMessageCount)
			c.Expect(msgcount, gs.Equals, int64(1))
			c.Expect(result, gs.Equals, string(matchBytes))

			close(inChan)
			wg.Wait() // wait for output to finish shutting down
		})

		c.Specify("far end not initially listening", func() {
			inChanCall := oth.MockOutputRunner.EXPECT().InChan().AnyTimes()
			inChanCall.Return(inChan)
			oth.MockOutputRunner.EXPECT().LogError(gomock.Any()).AnyTimes()

			err := tcpOutput.Init(config)
			c.Assume(err, gs.IsNil)

			pack.Message.SetPayload(outStr)
			go func() {
				wg.Add(1)
				c.Expect(err, gs.IsNil)
				err = tcpOutput.Run(oth.MockOutputRunner, oth.MockHelper)
				wg.Done()
			}()
			msgcount := atomic.LoadInt64(&tcpOutput.processMessageCount)
			c.Expect(msgcount, gs.Equals, int64(0))

			inChan <- pack

			for x := 0; x < 5 && msgcount == 0; x++ {
				msgcount = atomic.LoadInt64(&tcpOutput.processMessageCount)
				time.Sleep(time.Duration(100) * time.Millisecond)
			}

			// After the message is queued start the collector.
			// However, we don't have a way guarantee a send attempt has already
			// been made and that we are actually exercising the retry code.
			collectData := func(ch chan string) {
				ln, err := net.Listen("tcp", "localhost:9125")
				if err != nil {
					ch <- err.Error()
					return
				}
				conn, err := ln.Accept()
				if err != nil {
					ch <- err.Error()
					return
				}
				b := make([]byte, 1000)
				n, _ := conn.Read(b)
				ch <- string(b[0:n])
				conn.Close()
				ln.Close()
			}
			ch := make(chan string, 1) // don't block on put
			go collectData(ch)
			result := <-ch
			c.Expect(result, gs.Equals, string(matchBytes))

			close(inChan)
			wg.Wait() // wait for output to finish shutting down
		})
	})
}
