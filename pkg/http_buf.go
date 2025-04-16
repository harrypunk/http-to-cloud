package feature

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpBuf struct {
	Size int
}

func (hb *HttpBuf) Get(ctx context.Context, url string, ch chan<- []byte) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get error: %v", err)
	}

	chunkSize := 1024 * 1024 * hb.Size

	buffer := make([]byte, chunkSize)

	for {
		// Read a chunk from the response body
		n, err := io.ReadFull(resp.Body, buffer)
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
				return fmt.Errorf("context done") // Exit if context is cancelled
			}
		}

		// Check for errors after processing the read data
		if err != nil {
			if err.Error() == "EOF" || errors.Is(err, context.Canceled) || errors.Is(err, http.ErrBodyReadAfterClose) {
				// End of file or expected closure/cancellation
				log.Println("Finished reading response body.")
				break // Exit loop
			} else if err == io.ErrUnexpectedEOF {
				log.Printf("unexpected eof: %v", n)
				break
			} else {
				// Unexpected error
				log.Printf("Error reading response body: %v", err)
				break // Exit loop on error
			}
		}
	}

	log.Println("Download and piping finished.")
	close(ch)
	return nil
}
