package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Screenshot struct {
	Url        string
	Selector   string
	Quality    int
	ImageBytes *[]byte
	Viewport   page.Viewport
}

func (s *Screenshot) TakeScreenshot() error {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	defer cancel()

	viewport := page.Viewport{Width: 1025, Height: 820}

	if s.Viewport.Width != 0.0 {
		viewport = page.Viewport{Width: s.Viewport.Width, Height: s.Viewport.Height}
	}

	err := chromedp.Run(ctx, chromedp.Tasks{
		// chromedp.Navigate(s.Url),
		chromedp.EmulateViewport(int64(viewport.Width), int64(viewport.Height)),
		enableLifeCycleEvents(),
		navigateAndWaitFor(s.Url, "networkIdle"),
		chromedp.Screenshot(s.Selector, s.ImageBytes, chromedp.NodeVisible),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Screenshot) TakeFullScreenshot() error {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	defer cancel()

	err := chromedp.Run(ctx, chromedp.Tasks{
		// chromedp.Navigate(s.Url),
		enableLifeCycleEvents(),
		navigateAndWaitFor(s.Url, "networkIdle"),
		chromedp.FullScreenshot(s.ImageBytes, s.Quality),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}
		err = page.SetLifecycleEventsEnabled(true).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		chromedp.Navigate(url).Do(ctx)
		return waitFor(ctx, eventName)
	}
}

// waitFor blocks until eventName is received.
// Examples of events you can wait for:
//
//	init, DOMContentLoaded, firstPaint,
//	firstContentfulPaint, firstImagePaint,
//	firstMeaningfulPaintCandidate,
//	load, networkAlmostIdle, firstMeaningfulPaint, networkIdle
func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func GenerateImageFromTemplate(temp map[string]interface{}) string {

	reqBody, err := json.Marshal(temp)
	if err != nil {
		log.Printf("Error mashaling request body:: %s ", err.Error())
		return ""
	}

	req, err := http.NewRequest("POST", os.Getenv("HTCI_BASE"), bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("unable to create new request: %s", err.Error())
		return ""
	}
	username := os.Getenv("HTCI_USERNAME")
	password := os.Getenv("HTCI_PASSWORD")
	req.SetBasicAuth(username, password)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request was unsuccessful: %s", err.Error())
		return ""
	}

	defer resp.Body.Close()

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("unable to read response body: %s", err.Error())
		return ""
	}

	return result["url"]
}

func ConvertImage2Byte(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
