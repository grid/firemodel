// DO NOT EDIT - Code generated by firemodel v0.0.34-2-g6dab652-dirty.

package firemodel

import firestore "cloud.google.com/go/firestore"

type Client struct {
	Client *firestore.Client
}

func NewClient(client *firestore.Client) *Client {
	temp := &Client{Client: client}
	return temp
}
