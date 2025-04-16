package feature

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HttpBuf struct {
	Size int
}

func (hb *HttpBuf) Get(ctx context.Context, url string) (chan<- []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %v", err)
	}

	chunkSize := 1024 * 1024 * hb.Size

	bufReader := bufio.NewReaderSize(resp.Body, chunkSize)
	buffer := make([]byte, chunkSize)

	ch := make(chan []byte)

	for {
		// Read a chunk from the response body
		n, err := bufReader.Read(buffer)
		if n > 0 {
			// Important: Send a *copy* of the relevant part of the buffer.
			// The buffer will be reused in the next iteration.
			chunkCopy := make([]byte, n)
			copy(chunkCopy, buffer[:n])

			// Send the chunk copy over the channel
			select {
			case ch <- chunkCopy:
				// Chunk sent successfully
			case <-ctx.Done():
				log.Println("Download cancelled by context.")
				return nil, fmt.Errorf("context done") // Exit if context is cancelled
			}
		}

		// Check for errors after processing the read data
		if err != nil {
			if err == bufio.ErrBufferFull {
				// This shouldn't happen with Read if the buffer size matches chunkSize,
				// but handle defensively. Usually indicates buffer smaller than needed.
				log.Println("Warning: Buffer full during read, continuing.")
				continue // Try reading again
			} else if err.Error() == "EOF" || errors.Is(err, context.Canceled) || errors.Is(err, http.ErrBodyReadAfterClose) {
				// End of file or expected closure/cancellation
				log.Println("Finished reading response body.")
				break // Exit loop
			} else {
				// Unexpected error
				log.Printf("Error reading response body: %v", err)
				break // Exit loop on error
			}
		}
	}

	log.Println("Download and piping finished.")
	return ch, nil
}
