package easi

import(
	"context"
	"io"
	"time"
	"bytes"
	"encoding/csv"
	"github.com/jszwec/csvutil"
)

type Standard997V2 struct {
	EnvelopeHeaderV3 EnvelopeHeaderV3
	Body Standard997V2Body
	EnvelopeTrailerV3 EnvelopeTrailerV3
}

type Standard997V2Body struct {
	Header string
	TransactionType string
	VersionNumber string
	SenderQualifier string
	SenderID string
	ReceiverQualifier string
	ReceiverID string
	FileCreationDate string
	FileCreationTime string
	ProductionOrTest string
	InterchangeID string
	TransactionSetAcknowledgementCodes string
}

func (s *Standard997V2) Prep(ctx context.Context) (error){
	
	// Header
	errHeader := s.EnvelopeHeaderV3.Prep(ctx)
	if errHeader != nil {
		return errHeader
	}
	s.EnvelopeHeaderV3.TransactionType = "997"

	// Transaction
	s.Body.Header = "01"
	s.Body.TransactionType = "997"
	s.Body.VersionNumber = "2.0"
	s.Body.SenderQualifier = "01"
	s.Body.FileCreationDate = time.Now().Format("20060102")
	s.Body.FileCreationTime = time.Now().Format("150405")

	// Trailer
	errTrailer := s.EnvelopeTrailerV3.Prep(ctx)
	if errTrailer != nil {
		return errTrailer
	}

	return nil
}

func (s *Standard997V2) ToBytes(ctx context.Context) (*[]byte, error){

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
    w.Comma = '\t'
	enc := csvutil.NewEncoder(w)

	// Prep
	errPrep := s.Prep(ctx)
	if errPrep != nil {
		return nil, errPrep
	}
	
	// Headerless
	type header struct{}
	errHeader := enc.EncodeHeader(header{})
	if errHeader != nil {
		return nil, errHeader
	}

	// Envelope Header
	errEnvelopeHeaderV3 := enc.Encode(s.EnvelopeHeaderV3)
	if errEnvelopeHeaderV3 != nil {
		return nil, errEnvelopeHeaderV3
	}
	
	// Body
	errBody := enc.Encode(s.Body)
	if errBody != nil {
		return nil, errBody
	}

	// Envelope Trailer
	errEnvelopeTrailerV3 := enc.Encode(s.EnvelopeTrailerV3)
	if errEnvelopeTrailerV3 != nil {
		return nil, errEnvelopeTrailerV3
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	byteArray := buf.Bytes()

	return &byteArray, nil
}

func (s *Standard997V2) FromBytes(ctx context.Context, req []byte) (error){

	r := csv.NewReader(bytes.NewReader(req))
	r.Comma = '\t'

	// Headerless
	blankHeader, errHeader := csvutil.Header(Header{}, "csv")
	if errHeader != nil {
		return errHeader
	}

	// Decoder
	dec, errDecoder := csvutil.NewDecoder(r, blankHeader...)
	if errDecoder != nil {
		return errDecoder
	}

	for {
		var v struct{}
		if err := dec.Decode(&v); err == io.EOF {
			break
		}

		// Record
		var lineType string
		record := dec.Record()
		if len(record) > 0 {
			lineType = record[0]
		}

		// Build
		switch lineType {
		case "EASI":
			var x EnvelopeHeaderV3
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeHeaderV3 = x
		case "01":
			var x Standard997V2Body
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.Body = x
		case "EASX":
			var x EnvelopeTrailerV3
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s.EnvelopeTrailerV3 = x
		default:
			
		}

	}
	
	return nil
}

func (s *Standard997V2Body) FromSlice(ctx context.Context, req []string) (error){

	if len(req) > 0 {
		s.Header = req[0]
	}
	if len(req) > 1 {
		s.TransactionType = req[1]
	}
	if len(req) > 2 {
		s.VersionNumber = req[2]
	}
	if len(req) > 3 {
		s.SenderQualifier = req[3]
	}
	if len(req) > 4 {
		s.SenderID = req[4]
	}
	if len(req) > 5 {
		s.ReceiverQualifier = req[5]
	}
	if len(req) > 6 {
		s.ReceiverID = req[6]
	}
	if len(req) > 7 {
		s.FileCreationDate = req[7]
	}
	if len(req) > 8 {
		s.FileCreationTime = req[8]
	}
	if len(req) > 9 {
		s.ProductionOrTest = req[9]
	}
	if len(req) > 10 {
		s.InterchangeID = req[10]
	}
	if len(req) > 11 {
		s.TransactionSetAcknowledgementCodes = req[11]
	}
	return nil
}