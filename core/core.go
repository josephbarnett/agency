package core

import (
	"context"
	"fmt"
)

// Pipe is basic building block. Pipes can be composed together into pipeline
type Pipe func(context.Context, Message) (Message, error)

// Then takes a `next` pipe and returns new pipe that wraps `next`
func (p Pipe) Then(next Pipe) Pipe {
	return func(ctx context.Context, bb Message) (Message, error) {
		bb, err := p(ctx, bb)
		if err != nil {
			return nil, err
		}
		return next(ctx, bb)
	}
}

// Execute executes the pipe(line). This is syntactic sugar of regular function call
func (p Pipe) Execute(ctx context.Context, bb Message) (Message, error) {
	return p(ctx, bb)
}

// Message represents abstract message
type Message interface {
	Bytes() []byte
}

type TextMessage struct {
	Role    Role
	Content string
}

func (t TextMessage) Bytes() []byte {
	return []byte(t.Content)
}

// Bind allows to use prompt as a template by replacing printf directives like `%s` with the given `args`
func (t TextMessage) Bind(args ...any) TextMessage {
	return TextMessage{
		Role:    t.Role,
		Content: fmt.Sprintf(t.Content, args...),
	}
}

type ImageMessage struct {
	bb []byte
}

func (i ImageMessage) Bytes() []byte {
	return i.bb
}

func NewImageMessage(bb []byte) ImageMessage {
	return ImageMessage{bb}
}

type Role string

const (
	UserRole      Role = "user"
	SystemRole    Role = "system"
	AssistantRole Role = "assistant"
)

// NewUserMessage creates new `TextMessage` with the `Role` equal to `user`
func NewUserMessage(content string) TextMessage {
	return TextMessage{Role: UserRole, Content: content}
}

// NewSystemMessage creates new `TextMessage` with the `Role` equal to `system`
func NewSystemMessage(content string) TextMessage {
	return TextMessage{Role: SystemRole, Content: content}
}

type SpeechMessage struct {
	bb []byte
}

func (s SpeechMessage) Bytes() []byte {
	return s.bb
}

func NewSpeechMessage(bb []byte) SpeechMessage {
	return SpeechMessage{
		bb: bb,
	}
}
