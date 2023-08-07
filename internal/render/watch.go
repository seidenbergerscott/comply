package render

import (
	"net/http"
	"time"

	"github.com/gohugoio/hugo/watcher"
)

func watch(errCh chan error) {
	// TODO: study about the poll duration
	b, err := watcher.New(300 * time.Millisecond, 0, false)
	if err != nil {
		errCh <- err
		return
	}
	b.Add("./templates/")
	b.Add("./narratives/")
	b.Add("./architectures/")
	b.Add("./policies/")
	b.Add("./procedures/")

	b.Add("./.comply/")
	b.Add("./.comply/cache")
	b.Add("./.comply/cache/tickets")

	go func() {
		for {
			select {
			case e := <-b.Errors():
				errCh <- e
			case <-b.Events:
				broadcast()
			}
		}
	}()

	serveWs := func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			errCh <- err
			return
		}
		<-subscribe()
		time.Sleep(500 * time.Millisecond)
		ws.Close()
	}

	http.HandleFunc("/ws", serveWs)

	return
}
