package openai

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TIP
// headers
// x-ratelimit-limit-requests
// x-ratelimit-limit-tokens
// x-ratelimit-remaining-requests
// x-ratelimit-remaining-tokens
// x-ratelimit-reset-requests
// x-ratelimit-reset-tokens

type CompletionRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
}
type CompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Refusal     interface{}   `json:"refusal"`
			Annotations []interface{} `json:"annotations"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
			AudioTokens  int `json:"audio_tokens"`
		} `json:"prompt_tokens_details"`
		CompletionTokensDetails struct {
			ReasoningTokens          int `json:"reasoning_tokens"`
			AudioTokens              int `json:"audio_tokens"`
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	ServiceTier       string `json:"service_tier"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type OpenAI struct {
}

func (o *OpenAI) ParseRequest(req *http.Request) error {

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, req.Body)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(buf)

	var request CompletionRequest
	if err := json.Unmarshal(buf.Bytes(), &request); err != nil {
		return err
	}

	fmt.Println(request)
	return nil
}

func (o *OpenAI) ParseResponse(resp *http.Response, writer http.ResponseWriter) error {

	isStreaming := strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream")
	if isStreaming {
		if err := o.StreamParser(resp, writer); err != nil {
			return err
		}
		return nil
	}

	respBuf := &bytes.Buffer{}
	if _, err := io.Copy(respBuf, resp.Body); err != nil {
		return err
	}

	copyBody := respBuf.Bytes()
	_, err := io.Copy(writer, io.NopCloser(respBuf))
	if err != nil {
		return err
	}

	encoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(encoding, "gzip") {
		gzipReader, err := gzip.NewReader(bytes.NewReader(copyBody))
		if err != nil {
			return err
		}

		defer gzipReader.Close()
		uncompressedBuf := &bytes.Buffer{}
		if _, err = io.Copy(uncompressedBuf, gzipReader); err != nil {
			return err
		} else {
			// 압축 해제된 내용 출력
			fmt.Println(uncompressedBuf.String())
		}
	}

	return nil
}

func (o *OpenAI) StreamParser(resp *http.Response, writer http.ResponseWriter) error {
	reader := bufio.NewReader(resp.Body)
	var fullResponse strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		// 응답 캡처
		fullResponse.Write(line)

		// 클라이언트에 데이터 전송
		_, err = writer.Write(line)
		if err != nil {
			return err
		}

		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	return nil
}
