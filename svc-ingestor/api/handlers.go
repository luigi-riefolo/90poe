package api

import (
	"context"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/luigi-riefolo/nlp/svc-ingestor/pb"
	storerpb "github.com/luigi-riefolo/nlp/svc-storer/pb"
	"github.com/pkg/errors"
	"github.com/ttacon/libphonenumber"
)

// IngestFileHandler processes ingestion requests.
func (i *IngestorService) IngestFileHandler(ctx context.Context, req *pb.IngestRequest) (*pb.IngestResponse, error) {

	log.Debug("IngestFileHandler")

	return &pb.IngestResponse{}, nil
}

// process the data file and send the entries to the storer service.
// NOTE: this function can be triggered via an endpoint.
// Furthermore any malformed or invalid entry is skipped.
// The developer would consider the implementation of a list of failed user
// entries, so that the list can be reported and/or retried.
func (i *IngestorService) ingest() error {

	log.WithField("data_file", i.config.DataFile).Info("ingest")

	dataFile, err := os.Open(i.config.DataFile)
	if err != nil {
		return errors.Wrap(err, "could not open the data file")
	}
	defer dataFile.Close()

	// the channel that receives the parsed entry
	chn := make(chan storerpb.Entry)

	ctx, cancel := context.WithTimeout(context.Background(), storeClientTimeout)

	// process each entry
	go func(chn <-chan storerpb.Entry) {
		for entry := range chn {

			log.Debugf("sending to store service entry '%s'", entry.Id)

			req := pool.Get().(*storerpb.StoreEntryRequest)

			// normalise the phone number
			entry.MobileNumber, err = formatPhoneNumber(entry.MobileNumber)
			if err != nil {
				log.WithError(err).
					WithField("id", entry.Id).
					Error("could not format phone number")
				pool.Put(req)

				continue
			}

			req.Entry = &entry

			if _, err := i.storerClient.StoreEntryHandler(ctx, req); err != nil {
				log.WithError(err).Error("store client error")
			}

			pool.Put(req)
		}
		log.Info("file ingestion terminated")
		defer cancel()
	}(chn)

	if err := gocsv.UnmarshalToChan(dataFile, chn); err != nil {
		log.Fatal(err)
	}

	return nil
}

// formatPhoneNumber returns a formatted
// phone number that complies with RFC3966.
func formatPhoneNumber(num string) (string, error) {
	parsedNum, err := libphonenumber.Parse(num, phoneRegion)
	if err != nil {
		return "", err
	}
	return libphonenumber.Format(parsedNum, libphonenumber.RFC3966), nil
}
