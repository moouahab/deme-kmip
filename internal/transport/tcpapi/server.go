package tcpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/transport/kmipwire"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase"
)

type Response struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Server struct {
	dispatcher *usecase.Dispatcher
	collector  metrics.Collector
}

func NewServer(dispatcher *usecase.Dispatcher, collector metrics.Collector) *Server {
	return &Server{
		dispatcher: dispatcher,
		collector:  collector,
	}
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcpapi: listen: %w", err)
	}
	defer listener.Close()

	log.Printf("starting kmip demo tcp transport on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return fmt.Errorf("tcpapi: accept: %w", err)
			}
		}

		go func() {
			defer conn.Close()
			if err := s.HandleMessage(ctx, conn, conn); err != nil {
				log.Printf("tcpapi: failed to handle connection: %v", err)
			}
		}()
	}
}

func (s *Server) HandleMessage(ctx context.Context, reader io.Reader, writer io.Writer) error {
	s.collector.IncHTTPRequests(ctx)

	data, err := readLimited(reader, ttlv.MaxValueLength)
	if err != nil {
		s.collector.IncHTTPErrors(ctx)
		return writeTCPResponse(writer, Response{
			OK:      false,
			Error:   "request_too_large",
			Message: err.Error(),
		})
	}

	blocks, err := ttlv.DecodeBlocks(data)
	if err != nil {
		s.collector.IncHTTPErrors(ctx)
		return writeTCPResponse(writer, Response{
			OK:      false,
			Error:   "bad_request",
			Message: err.Error(),
		})
	}

	req, err := kmipwire.BlocksToOperationRequest(blocks)
	if err != nil {
		s.collector.IncHTTPErrors(ctx)
		return writeTCPResponse(writer, Response{
			OK:      false,
			Error:   "bad_request",
			Message: err.Error(),
		})
	}

	switch req.Operation {
	case ttlv.OperationCreate:
		s.collector.IncCreateKey(ctx)
	case ttlv.OperationGet:
		s.collector.IncGetKey(ctx)
	case ttlv.OperationDestroy:
		s.collector.IncDestroyKey(ctx)
	}

	resp, err := s.dispatcher.Dispatch(ctx, req)
	if err != nil {
		if errors.Is(err, kms.ErrKeyNotFound) {
			s.collector.IncNotFound(ctx)
			return writeTCPResponse(writer, Response{
				OK:      false,
				Error:   "not_found",
				Message: err.Error(),
			})
		}

		s.collector.IncHTTPErrors(ctx)
		return writeTCPResponse(writer, Response{
			OK:      false,
			Error:   "bad_request",
			Message: err.Error(),
		})
	}

	s.collector.IncSuccess(ctx)
	return writeTCPResponse(writer, Response{
		OK:   true,
		Data: resp,
	})
}

func readLimited(reader io.Reader, maxBytes int) ([]byte, error) {
	limited := io.LimitReader(reader, int64(maxBytes)+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("cannot read request: %w", err)
	}
	if len(data) > maxBytes {
		return nil, fmt.Errorf("request body too large")
	}
	return data, nil
}

func writeTCPResponse(writer io.Writer, resp Response) error {
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("tcpapi: encode response: %w", err)
	}
	return nil
}
