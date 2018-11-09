// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains functions for interacting with the Discord REST/JSON API
// at the lowest level.

package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/robbix1206/discordgo/logging"
)

// All error constants
var (
	ErrJSONUnmarshal = errors.New("json unmarshal")
	ErrUnauthorized  = errors.New("HTTP request was unauthorized. This could be because the provided token was not a bot token. Please add \"Bot \" to the start of your token. https://discordapp.com/developers/docs/reference#authentication-example-bot-token-authorization-header")
)

// RequestWithBucket is the same as RequestWithBucketID but the bucket id is the same as the urlStr
func (s *Client) RequestWithBucket(method, urlStr string, data interface{}) (response []byte, err error) {
	return s.RequestWithBucketID(method, urlStr, data, strings.SplitN(urlStr, "?", 2)[0])
}

// RequestWithBucketID makes a (GET/POST/...) Requests to Discord REST API with JSON data.
func (s *Client) RequestWithBucketID(method, urlStr string, data interface{}, bucketID string) (response []byte, err error) {
	var body []byte
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	return s.Request(method, urlStr, "application/json", body, bucketID, 0)
}

// Request makes a (GET/POST/...) Requests to Discord REST API.
// Sequence is the sequence number, if it fails with a 502 it will
// retry with sequence+1 until it either succeeds or sequence >= session.MaxRestRetries
func (s *Client) Request(method, urlStr, contentType string, b []byte, bucketID string, sequence int) (response []byte, err error) {
	if bucketID == "" {
		bucketID = strings.SplitN(urlStr, "?", 2)[0]
	}
	return s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucket(bucketID), sequence)
}

// RequestWithLockedBucket makes a request using a bucket that's already been locked
func (s *Client) RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *Bucket, sequence int) (response []byte, err error) {
	if s.LogLevel >= logging.LogDebug {
		log.Printf("API REQUEST %8s :: %s\n", method, urlStr)
		log.Printf("API REQUEST  PAYLOAD :: [%s]\n", string(b))
	}

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(b))
	if err != nil {
		bucket.Release(nil)
		return
	}

	// Not used on initial login..
	// TODO: Verify if a login, otherwise complain about no-token
	if s.Token != "" {
		req.Header.Set("authorization", s.Token)
	}

	req.Header.Set("Content-Type", contentType)
	// TODO: Make a configurable static variable.
	req.Header.Set("User-Agent", "DiscordBot (https://github.com/bwmarrin/discordgo, v"+version+")")

	if s.LogLevel >= logging.LogDebug {
		for k, v := range req.Header {
			log.Printf("API REQUEST   HEADER :: [%s] = %+v\n", k, v)
		}
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		bucket.Release(nil)
		return
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			log.Println("error closing resp body")
		}
	}()

	err = bucket.Release(resp.Header)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if s.LogLevel >= logging.LogDebug {

		log.Printf("API RESPONSE  STATUS :: %s\n", resp.Status)
		for k, v := range resp.Header {
			log.Printf("API RESPONSE  HEADER :: [%s] = %+v\n", k, v)
		}
		log.Printf("API RESPONSE    BODY :: [%s]\n\n\n", response)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
	case http.StatusBadGateway:
		// Retry sending request if possible
		if sequence < s.MaxRestRetries {

			s.log(logging.LogInformational, "%s Failed (%s), Retrying...", urlStr, resp.Status)
			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucketObject(bucket), sequence+1)
		} else {
			err = fmt.Errorf("Exceeded Max retries HTTP %s, %s", resp.Status, response)
		}
	case 429: // TOO MANY REQUESTS - Rate limiting
		rl := TooManyRequests{}
		err = json.Unmarshal(response, &rl)
		if err != nil {
			s.log(logging.LogError, "rate limit unmarshal error, %s", err)
			return
		}
		s.log(logging.LogInformational, "Rate Limiting %s, retry in %d", urlStr, rl.RetryAfter)
		// Is this event really useful ? Shouldn't be
		//s.handleEvent(rateLimitEventType, RateLimit{TooManyRequests: &rl, URL: urlStr})

		time.Sleep(rl.RetryAfter * time.Millisecond)
		// we can make the above smarter
		// this method can cause longer delays than required

		response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucketObject(bucket), sequence)
	case http.StatusUnauthorized:
		if strings.Index(s.Token, "Bot ") != 0 {
			s.log(logging.LogInformational, ErrUnauthorized.Error())
			err = ErrUnauthorized
		}
		fallthrough
	default: // Error condition
		err = newRestError(req, resp, response)
	}
	return
}

// Unmarshal unmarshal the given data
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return ErrJSONUnmarshal
	}

	return nil
}

func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

func (r RESTError) Error() string {
	return "HTTP " + r.Response.Status + ", " + string(r.ResponseBody)
}
