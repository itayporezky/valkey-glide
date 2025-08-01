// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package glide

// #include "lib.h"
import "C"

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/itayporezky/valkey-glide/go/v4/config"
	"github.com/itayporezky/valkey-glide/go/v4/models"
)

var (
	ErrPubSubPushInvalid       = errors.New("received invalid push: empty or in incorrect format")
	ErrPubSubPushMissingKind   = errors.New("received invalid push: missing kind field")
	ErrPubSubPushMissingValues = errors.New("received invalid push: missing values field")
)

type MessageCallbackError struct {
	cause error
}

func (e *MessageCallbackError) Error() string {
	return fmt.Sprintf("error in message callback: %v", e.cause)
}

func (e *MessageCallbackError) Cause() error {
	return e.cause
}

// *** Message Handler ***

type MessageHandler struct {
	callback config.MessageCallback
	context  any
	queue    *PubSubMessageQueue
}

func NewMessageHandler(callback config.MessageCallback, context any) *MessageHandler {
	return &MessageHandler{
		callback: callback,
		context:  context,
		queue:    NewPubSubMessageQueue(),
	}
}

func (handler *MessageHandler) handleMessage(message *models.PubSubMessage) error {
	if handler.callback != nil {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("panic in message callback", err.Error())
			}
		}()

		handler.callback(message, handler.context)
		return nil
	} else {
		handler.queue.Push(message)
		return nil
	}
}

func (handler *MessageHandler) GetQueue() *PubSubMessageQueue {
	return handler.queue
}

// *** Message Queue ***

type PubSubMessageQueue struct {
	mu                      sync.Mutex
	messages                []*models.PubSubMessage
	waiters                 []chan *models.PubSubMessage
	nextMessageReadyCh      chan struct{}
	nextMessageReadySignals []chan struct{}
}

func NewPubSubMessageQueue() *PubSubMessageQueue {
	return &PubSubMessageQueue{
		messages:                make([]*models.PubSubMessage, 0),
		waiters:                 make([]chan *models.PubSubMessage, 0),
		nextMessageReadyCh:      make(chan struct{}, 1),
		nextMessageReadySignals: make([]chan struct{}, 0),
	}
}

func (queue *PubSubMessageQueue) Push(message *models.PubSubMessage) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	// If there's a waiter, deliver the message directly
	if len(queue.waiters) > 0 {
		waiterCh := queue.waiters[0]
		queue.waiters = queue.waiters[1:]
		waiterCh <- message
		return
	}

	// Otherwise, add to the queue
	queue.messages = append(queue.messages, message)

	// Signal that a new message is ready
	select {
	case queue.nextMessageReadyCh <- struct{}{}:
	default:
		// Channel already has a signal
	}

	// Signal any waiters
	for _, ch := range queue.nextMessageReadySignals {
		select {
		case ch <- struct{}{}:
		default:
			// Channel is full, receiver might not be listening
		}
	}
}

func (queue *PubSubMessageQueue) Pop() *models.PubSubMessage {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	if len(queue.messages) == 0 {
		return nil
	}

	message := queue.messages[0]
	queue.messages = queue.messages[1:]
	return message
}

func (queue *PubSubMessageQueue) WaitForMessage() <-chan *models.PubSubMessage {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	// If a message is already queued, return it immediately
	if len(queue.messages) > 0 {
		messageCh := make(chan *models.PubSubMessage, 1)
		message := queue.messages[0]
		queue.messages = queue.messages[1:]
		messageCh <- message
		return messageCh
	}

	// Otherwise register a waiter
	messageCh := make(chan *models.PubSubMessage, 1)
	queue.waiters = append(queue.waiters, messageCh)
	return messageCh
}

func (queue *PubSubMessageQueue) RegisterSignalChannel(ch chan struct{}) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.nextMessageReadySignals = append(queue.nextMessageReadySignals, ch)

	// Add this to signal if there are existing messages
	if len(queue.messages) > 0 {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (queue *PubSubMessageQueue) UnregisterSignalChannel(ch chan struct{}) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	for idx, channel := range queue.nextMessageReadySignals {
		if channel == ch {
			queue.nextMessageReadySignals = append(
				queue.nextMessageReadySignals[:idx],
				queue.nextMessageReadySignals[idx+1:]...)
			break
		}
	}
}
