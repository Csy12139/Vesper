package common

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
)

// Timeout for MN client operations
const mnClientTimeout = 5 * time.Second

// MNClient represents a client connection to the MN server
type MNClient struct {
	addr   string
	conn   *grpc.ClientConn
	client pb.DN2MNClient
	mu     sync.Mutex
}

// NewMNClient creates a new MN client
func NewMNClient(addr string) *MNClient {
	return &MNClient{
		addr: addr,
	}
}

// getClient returns the gRPC client, creating a new connection if needed
func (c *MNClient) getClient() (pb.DN2MNClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MN: %w", err)
		}
		c.conn = conn
		c.client = pb.NewDN2MNClient(conn)
	}
	return c.client, nil
}

// Close closes the client connection
func (c *MNClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// DoHeartbeat sends a heartbeat to the MN server and returns the response
func (c *MNClient) DoHeartbeat(uuid string) (*HeartbeatResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &HeartbeatRequest{UUID: uuid}
	pbReq := HeartbeatRequest2Proto(req)
	pbResp, err := client.DoHeartbeat(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("heartbeat failed: %w", err)
	}

	return Proto2HeartbeatResponse(pbResp), nil
}
