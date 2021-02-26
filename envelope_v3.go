package easi

import(
	"context"
	"time"
	"strconv"
)

type EnvelopeHeaderV3 struct {
	Header string
	VersionNumber string
	SenderQualifier string
	SenderID string
	ReceiverQualifier string
	ReceiverID string
	FileCreationDate string
	FileCreationTime string
	TimeZone string
	ProductionOrTest string
	TransactionType string
	InterchangeID string
}

type EnvelopeTrailerV3 struct {
	RoutingTrailerRecord string
	InterchangeID string
	NumberOfDocuments int
}

func (s *EnvelopeHeaderV3) Prep(ctx context.Context) (error){

	s.Header = "EASI"
	s.VersionNumber = "3.0"
	s.SenderQualifier = "01"
	s.ReceiverQualifier = "01"
	s.FileCreationDate = time.Now().Format("20060102")
	s.FileCreationTime = time.Now().Format("150405")
	s.TimeZone = "UTC"
	if s.ProductionOrTest == "" {
		s.ProductionOrTest = "T"
	}

	return nil
}

func (s *EnvelopeTrailerV3) Prep(ctx context.Context) (error){

	s.RoutingTrailerRecord = "EASX"
	s.NumberOfDocuments = 1

	return nil
}

func (s *EnvelopeHeaderV3) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.Header = req[0]
	}
	if len(req) > 1 {
		s.VersionNumber = req[1]
	}
	if len(req) > 2 {
		s.SenderQualifier = req[2]
	}
	if len(req) > 3 {
		s.SenderID = req[3]
	}
	if len(req) > 4 {
		s.ReceiverQualifier = req[4]
	}
	if len(req) > 5 {
		s.ReceiverID = req[5]
	}
	if len(req) > 6 {
		s.FileCreationDate = req[6]
	}
	if len(req) > 7 {
		s.FileCreationTime = req[7]
	}
	if len(req) > 8 {
		s.TimeZone = req[8]
	}
	if len(req) > 9 {
		s.ProductionOrTest = req[9]
	}
	if len(req) > 10 {
		s.TransactionType = req[10]
	}
	if len(req) > 11 {
		s.InterchangeID = req[11]
	}

	return nil
}

func (s *EnvelopeTrailerV3) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.RoutingTrailerRecord = req[0]
	}
	if len(req) > 1 {
		s.InterchangeID = req[1]
	}
	if len(req) > 2 {
		if req[2] != "" {
			numberOfDocuments, err := strconv.Atoi(req[2])
			if err != nil {
				return err
			}
			s.NumberOfDocuments = numberOfDocuments
		}
	}
	

	return nil
}