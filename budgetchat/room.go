package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type room struct {
	users map[string]net.Conn
	sync.RWMutex
}

func NewRoom() *room {
	return &room{
		users: make(map[string]net.Conn),
	}
}

func (r *room) addUser(userName string, conn net.Conn) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.users[userName]; ok {
		return errors.New("username already exists")
	}

	r.users[userName] = conn

	err := r.publish([]byte(fmt.Sprintf("* %s has entered the room\n", userName)), userName)

	return err
}

func (r *room) removeUser(userName string, conn net.Conn) error {
	r.Lock()
	defer r.Unlock()

	delete(r.users, userName)
	err := r.publish([]byte(fmt.Sprintf("* %s has left the room\n", userName)), userName)

	return err
}

func (r *room) publish(msg []byte, sender string) error {
	for user, conn := range r.users {
		if sender == user {
			continue
		}

		conn.Write(msg)
	}

	return nil
}

func (r *room) getUsers() []string {
	r.RLock()
	defer r.RUnlock()
	var users []string

	for user := range r.users {
		users = append(users, user)
	}

	return users
}
