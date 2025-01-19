package libvirt

import (
	"libvirt.org/go/libvirt"
)

type Client struct {
	conn *libvirt.Connect
}

func NewClient() (*Client, error) {
	conn, err := libvirt.NewConnectWithAuth("qemu:///system", &libvirt.ConnectAuth{}, libvirt.ConnectFlags(0))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		_, err := c.conn.Close()
		return err
	}
	return nil
}

func (c *Client) GetConnection() *libvirt.Connect {
	return c.conn
}

func (c *Client) IsConnected() (bool, error) {
	return c.conn.IsAlive()
}
