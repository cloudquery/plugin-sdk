package state

import "context"

type Client struct {
	// managedPlugin managedplugin.Client
}

// func NewState(ctx context.Context, managedPlugin managedplugin.Client) *Client {
// 	return &Client{
// 		managedPlugin: managedPlugin,
// 	}
// 	c := pbPlugin.NewPluginClient(managedPlugin.Conn)
// 	c.Write(ctx, )
// }

func NewState(spec any) *Client {
	return &Client{}
}

func (* Client) SetKey(ctx context.Context, key string, value string) error {
	return nil
}

func (* Client) GetKey(ctx context.Context, key string) (string, error) { 
	return "", nil
}