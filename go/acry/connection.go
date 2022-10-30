package acry

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type connection struct {
	conn net.Conn
	i    uint
}

type connectionPool struct {
	config      *acryConfig
	connections []*connection
	mux         *sync.Mutex
	poolsize    uint
	inuse       uint
	ctx         context.Context
	cancelCtx   context.CancelFunc
}

func Acry(ctx context.Context, config *acryConfig) *connectionPool {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &connectionPool{
		config:      config,
		connections: getConnectionPool(config),
		mux:         &sync.Mutex{},
		poolsize:    config.required,
		inuse:       0,
		ctx:         ctx,
		cancelCtx:   cancelFunc,
	}
}

func getConnectionPool(config *acryConfig) []*connection {
	var connections []*connection
	for i := uint(0); i < config.GetRequiredConnectionPool(); i++ {
		connections = append(connections, getConnection(config.host, config.GetRetryInterval(), i))
	}
	return connections
}

func getConnection(host string, in uint, index uint) *connection {
	for {
		con, err := net.Dial("tcp", host)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(1000000 * in))
			continue
		}
		connection := &connection{
			conn: con,
			i:    index,
		}
		return connection
	}
}

func (c *connectionPool) popConnectionFromPool() *connection {
	var connection *connection
	for {
		c.mux.Lock()
		if c.poolsize-c.inuse > 0 {
			connection = c.connections[0]
			c.connections = c.connections[1:]
			break
		} else if c.config.GetMaximumConnectionPool() > c.poolsize {
			connection = getConnection(c.config.host, c.config.GetRetryInterval(), c.poolsize)
			c.poolsize++
			break
		}
		c.mux.Unlock()
		time.Sleep(time.Duration(c.config.retryInterval * uint(time.Millisecond)))
		// log.Println("retry acry connection . . .")
	}
	c.inuse++
	c.mux.Unlock()
	return connection
}

func (c *connectionPool) pushConnectionFromPool(con *connection) {
	c.mux.Lock()
	c.connections = append(c.connections, con)
	c.inuse--
	c.mux.Unlock()
}

func (c *connectionPool) Query(r *request) (*response, error) {
	requestRaw, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	requestRaw = append(requestRaw, '\r')
	connection := c.popConnectionFromPool()
	connectionIndex := connection.i

QUERY_RETRY:
	select {
	case <-c.ctx.Done():
		connection.conn.Close()
		return nil, errors.New("query execution stopped because context ended")
	default:
	}
	_, err = connection.conn.Write(requestRaw)
	if err != nil {
		connection = getConnection(c.config.host, c.config.GetRetryInterval(), connectionIndex)
		goto QUERY_RETRY
	}
	reader := bufio.NewReader(connection.conn)
	responseRaw, err := reader.ReadBytes('\r')
	if err != nil {
		connection = getConnection(c.config.host, c.config.GetRetryInterval(), connectionIndex)
		goto QUERY_RETRY
	}
	select {
	case <-c.ctx.Done():
		connection.conn.Close()
		return nil, errors.New("query execution stopped because context ended")
	default:
	}
	c.pushConnectionFromPool(connection)

	response := response{}
	json.Unmarshal(responseRaw, &response)
	response.conIndex = connectionIndex
	return &response, nil
}

func (c *connectionPool) Close() {
	c.cancelCtx()
	for _, conn := range c.connections {
		conn.conn.Close()
	}
}
