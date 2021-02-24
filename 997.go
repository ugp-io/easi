package easi

import(
	"context"
	"io"
	"time"
	"bytes"
	"encoding/csv"
	"github.com/jszwec/csvutil"
)

type Standard997 struct {
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

func (s *Standard997) Prep(ctx context.Context) (error){

	// Transaction
	s.Header = "01"
	s.TransactionType = "997"
	s.VersionNumber = "2.0"
	s.SenderQualifier = "01"
	s.FileCreationDate = time.Now().Format("20060102")
	s.FileCreationTime = time.Now().Format("150405")


	return nil
}

func (s *Standard997) ToBytes(ctx context.Context) (*[]byte, error){

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

	// Transaction
	errTransaction := enc.Encode(s)
	if errTransaction != nil {
		return nil, errTransaction
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	byteArray := buf.Bytes()

	return &byteArray, nil
}

func (s *Standard997) FromBytes(ctx context.Context, req []byte) (error){

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
		case "01":
			var x Standard997
			err := x.FromSlice(ctx, record)
			if err != nil {
				return err
			}
			s = &x
		default:
			
		}

	}
	
	return nil
}

func (s *Standard997) FromSlice(ctx context.Context, req []string) (error){

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