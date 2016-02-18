package models

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// QualityReport is a representation of a calculation of an electronic
// clinical quality measure
type QualityReport struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	NPI             string        `bson:"npi,omitempty" json:"npi,omitempty"`
	CalculationTime time.Time     `bson:"calculation_time,omitempty" json:"calculationTime,omitempty"`
	Status          Status        `bson:"status,omitempty" json:"status,omitempty"`
	MeasureID       string        `bson:"measure_id,omitempty" json:"measureId,omitempty" validate:"nonzero"`
	SubID           string        `bson:"sub_id,omitempty" json:"subId,omitempty"`
	EffectiveDate   int32         `bson:"effective_date,omitempty" json:"effectiveDate,omitempty" validate:"nonzero"`
}

// FindQualityAndPopulateQualityReport will attempt to find a QualityReport in
// the query_cache based on the measure id, sub id and effective date passed in.
// If it finds the associated document in the database, it will return true
// and populate the other fields in the QualityReport that is passed in.
// Otherwise, it will return false, and the passed in QualityReport will remain
// unchanged.
func FindQualityAndPopulateQualityReport(db *mgo.Database, qr *QualityReport) (bool, error) {
	query := bson.M{"measure_id": qr.MeasureID, "effective_date": qr.EffectiveDate}
	if qr.SubID != "" {
		query["sub_id"] = qr.SubID
	}
	q := db.C("query_cache").Find(query)
	count, err := q.Count()
	if err != nil {
		return false, err
	}
	switch count {
	case 0:
		return false, nil
	case 1:
		err = q.One(qr)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.New("Found more than one quality report for this")
	}
	return true, nil
}

func FindOrCreateQualityReport(db *mgo.Database, qr *QualityReport) error {
	exists, err := FindQualityAndPopulateQualityReport(db, qr)
	if err != nil {
		return err
	}
	if !exists {
		qr.ID = bson.NewObjectId()
		err = db.C("query_cache").Insert(qr)
		if err != nil {
			return err
		}
	}
	return nil
}

type Status struct {
	State string   `bson:"state,omitempty" json:"state,omitempty"`
	Log   []string `bson:"log,omitempty" json:"log,omitempty"`
}

type QualityReportResult struct {
	PopulationIDs            PopulationIDs `bson:"population_ids,omitempty" json:"populationIds,omitempty"`
	InitialPatientPopulation int32         `bson:"IPP" json:"initialPatientPopulation"`
	Denominator              int32         `bson:"DENOM,omitempty" json:"denominator,omitempty"`
	Exception                int32         `bson:"DENEXCP,omitempty" json:"exception,omitempty"`
	Exclusion                int32         `bson:"DENEX,omitempty" json:"exclusion,omitempty"`
	Numerator                int32         `bson:"NUMER,omitempty" json:"numerator,omitempty"`
	AntiNumerator            int32         `bson:"antinumerator,omitempty" json:"antinumerator,omitempty"`
	MeasurePopulation        int32         `bson:"MSRPOPL,omitempty" json:"measurePopulation,omitempty"`
	Observation              float32       `bson:"OBSERV,omitempty" json:"Observation,omitempty"`
}

type PopulationIDs struct {
	InitialPatientPopulation string `bson:"IPP,omitempty" json:"initialPatientPopulation,omitempty"`
	Denominator              string `bson:"DENOM,omitempty" json:"denominator,omitempty"`
	Exception                string `bson:"DENEXCP,omitempty" json:"exception,omitempty"`
	Exclusion                string `bson:"DENEX,omitempty" json:"exclusion,omitempty"`
	Numerator                string `bson:"NUMER,omitempty" json:"numerator,omitempty"`
	MeasurePopulation        string `bson:"MSRPOPL,omitempty" json:"measurePopulation,omitempty"`
	Observation              string `bson:"OBSERV,omitempty" json:"observation,omitempty"`
}
