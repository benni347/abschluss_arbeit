package main

type Intent struct {
	ID      IntentID
	Details interface{}
}

type IntentID = string

const (
	ViewReplyWithID IntentID = "view-reply-with"
)

type ViewReplyWithIDDetails struct {
	NodeID string
}
