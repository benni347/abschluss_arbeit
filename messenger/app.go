package main

import (
	"context"

	utils "github.com/benni347/messengerutils"
	webrtc "github.com/pion/webrtc/v3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) CreatePeerConnection() *webrtc.PeerConnection {
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		utils.PrintError(err)
		panic(err)
	}

	return peerConnection
}

func (a *App) TransmitDataText(peerConnection *webrtc.PeerConnection, data string) {
	protcol := "tcp"
	order := true
	var maxReTransmission uint16 = 5000

	dataChannelInit := webrtc.DataChannelInit{
		Ordered:        &order,
		MaxRetransmits: &maxReTransmission,
		Protocol:       &protcol,
	}
	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("text", &dataChannelInit)
	if err != nil {
		utils.PrintError("During the creation of the data channel an error ocured", err)
		panic(err)
	}

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		// Send the data
		dataChannel.SendText(data)
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		println("Received message: " + string(msg.Data))
	})
}
