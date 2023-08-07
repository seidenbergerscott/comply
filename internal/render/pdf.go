package render

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/strongdm/comply/internal/model"
)

func pdf(output string, live bool, errCh chan error, wg *sync.WaitGroup) {
	var pdfWG sync.WaitGroup

	errOutputCh := make(chan error)

	for {
		_, data, err := loadWithStats()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to load data")
			return
		}

		policies, err := model.ReadPolicies()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to read policies")
			return
		}
		for _, policy := range policies {
			renderToFilesystem(&pdfWG, errOutputCh, data, policy, live)
			err = <-errOutputCh
			if err != nil {
				errCh <- err
				wg.Done()
				return
			}
		}

		narratives, err := model.ReadNarratives()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to read narratives")
			return
		}

		for _, narrative := range narratives {
			renderToFilesystem(&pdfWG, errOutputCh, data, narrative, live)
			err = <-errOutputCh
			if err != nil {
				errCh <- err
				wg.Done()
				return
			}
		}

		architectures, err := model.ReadArchitectures()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to read architectures")
			return
		}

		for _, architecture := range architectures {
			renderToFilesystem(&pdfWG, errOutputCh, data, architecture, live)
			err = <-errOutputCh
			if err != nil {
				errCh <- err
				wg.Done()
				return
			}
		}


		pdfWG.Wait()

		if !live {
			wg.Done()
			return
		}
		<-subscribe()
	}
}
