package main

import (
	"fmt"
	"strings"

	"github.com/khulnasoft/tracker/signatures/helpers"
	"github.com/khulnasoft/tracker/types/detect"
	"github.com/khulnasoft/tracker/types/protocol"
	"github.com/khulnasoft/tracker/types/trace"
)

//
// HOWTO: The way to trigger this test signature is to execute:
//
//        curl google.com
//
//        This will cause it trigger once and reset it status.

type e2eHTTPRequest struct {
	cb detect.SignatureHandler
}

func (sig *e2eHTTPRequest) Init(ctx detect.SignatureContext) error {
	sig.cb = ctx.Callback
	return nil
}

func (sig *e2eHTTPRequest) GetMetadata() (detect.SignatureMetadata, error) {
	return detect.SignatureMetadata{
		ID:          "HTTPRequest",
		EventName:   "HTTPRequest",
		Version:     "0.1.0",
		Name:        "Network HTTP Request Test",
		Description: "Network E2E Tests: HTTP Request",
		Tags:        []string{"e2e", "network"},
	}, nil
}

func (sig *e2eHTTPRequest) GetSelectedEvents() ([]detect.SignatureEventSelector, error) {
	return []detect.SignatureEventSelector{
		{Source: "tracker", Name: "net_packet_http_request"},
	}, nil
}

func (sig *e2eHTTPRequest) OnEvent(event protocol.Event) error {
	eventObj, ok := event.Payload.(trace.Event)
	if !ok {
		return fmt.Errorf("failed to cast event's payload")
	}

	if eventObj.ProcessName != "curl" {
		return nil
	}

	if eventObj.EventName == "net_packet_http_request" {
		// validate tast context
		if eventObj.HostName == "" {
			return nil
		}

		httpRequest, err := helpers.GetProtoHTTPRequestByName(eventObj, "http_request")
		if err != nil {
			return err
		}

		if !strings.HasPrefix(httpRequest.Protocol, "HTTP/") {
			return nil
		}

		m, _ := sig.GetMetadata()
		sig.cb(&detect.Finding{
			SigMetadata: m,
			Event:       event,
			Data:        map[string]interface{}{},
		})
	}

	return nil
}

func (sig *e2eHTTPRequest) OnSignal(s detect.Signal) error {
	return nil
}

func (sig *e2eHTTPRequest) Close() {}
