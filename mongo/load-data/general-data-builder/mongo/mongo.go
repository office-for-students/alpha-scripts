package mongo

import (
	"errors"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// Mongo ...
type Mongo struct {
	URI     string
	Session *mgo.Session
}

// Init creates a new mgo.Session with a strong consistency and a write mode of "majortiy".
func (m *Mongo) Init() (session *mgo.Session, err error) {
	if session != nil {
		return nil, errors.New("session already exists")
	}

	if session, err = mgo.Dial(m.URI); err != nil {
		return nil, err
	}

	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	session.SetMode(mgo.Strong, true)
	return session, nil
}

// AddCAHCode ...
func (m *Mongo) AddCAHCode(database, collection string, cahObject *data.SubjectObject) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(cahObject); err != nil {
		log.ErrorC("failed to create CAH code data resource", err, nil)
		return
	}

	return
}

// AddCommonData ...
func (m *Mongo) AddCommonData(database, collection string, commonData *data.CommonData) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(commonData); err != nil {
		log.ErrorC("failed to create common data resource", err, nil)
		return
	}

	return
}

// AddContinuation ...
func (m *Mongo) AddContinuation(database, collection string, continuation *data.Continuation) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(continuation); err != nil {
		log.ErrorC("failed to create continuation resource", err, nil)
		return
	}

	return
}

// AddCourseLocation ...
func (m *Mongo) AddCourseLocation(database, collection string, courseLocation *data.Location) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(courseLocation); err != nil {
		log.ErrorC("failed to create course location resource", err, nil)
		return
	}

	return
}

// AddDegreeClass ...
func (m *Mongo) AddDegreeClass(database, collection string, degreeClass *data.DegreeClass) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(degreeClass); err != nil {
		log.ErrorC("failed to create degree class resource", err, nil)
		return
	}

	return
}

// AddEmployment ...
func (m *Mongo) AddEmployment(database, collection string, employment *data.Employment) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(employment); err != nil {
		log.ErrorC("failed to create employment resource", err, nil)
		return
	}

	return
}

// AddEntry ...
func (m *Mongo) AddEntry(database, collection string, entry *data.Entry) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(entry); err != nil {
		log.ErrorC("failed to create entry resource", err, nil)
		return
	}

	return
}

// AddInstitution ...
func (m *Mongo) AddInstitution(database, collection string, institution *data.Institution) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(institution); err != nil {
		log.ErrorC("failed to create raw institution resource", err, nil)
		return
	}

	return
}

// AddInstitutionLocation ...
func (m *Mongo) AddInstitutionLocation(database, collection string, institutionLocation *data.InstitutionLocation) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(institutionLocation); err != nil {
		log.ErrorC("failed to create institution location resource", err, nil)
		return
	}

	return
}

// AddJobList ...
func (m *Mongo) AddJobList(database, collection string, jobList *data.JobList) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(jobList); err != nil {
		log.ErrorC("failed to create job list resource", err, nil)
		return
	}

	return
}

// AddJobType ...
func (m *Mongo) AddJobType(database, collection string, jobType *data.JobType) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(jobType); err != nil {
		log.ErrorC("failed to create job typr resource", err, nil)
		return
	}

	return
}

// AddLEOCourseStatistic ...
func (m *Mongo) AddLEOCourseStatistic(database, collection string, courseLocation *data.Leo) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(courseLocation); err != nil {
		log.ErrorC("failed to create leo course statistic resource", err, nil)
		return
	}

	return
}

// AddNHSNSS ...
func (m *Mongo) AddNHSNSS(database, collection string, nss *data.NHSNSS) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(nss); err != nil {
		log.ErrorC("failed to create nhs nss resource", err, nil)
		return
	}

	return
}

// AddNSS ...
func (m *Mongo) AddNSS(database, collection string, nss *data.NSS) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(nss); err != nil {
		log.ErrorC("failed to create nss resource", err, nil)
		return
	}

	return
}

// AddQualification ...
func (m *Mongo) AddQualification(database, collection string, qualification *data.Qualification) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(qualification); err != nil {
		log.ErrorC("failed to create qualification resource", err, nil)
		return
	}

	return
}

// AddSalary ...
func (m *Mongo) AddSalary(database, collection string, salary *data.Salary) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(salary); err != nil {
		log.ErrorC("failed to create salary resource", err, nil)
		return
	}

	return
}

// AddSubject ...
func (m *Mongo) AddSubject(database, collection string, subject *data.Subject) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(subject); err != nil {
		log.ErrorC("failed to create subject resource", err, nil)
		return
	}

	return
}

// AddTariff ...
func (m *Mongo) AddTariff(database, collection string, tariff *data.Tariff) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(tariff); err != nil {
		log.ErrorC("failed to create tariff resource", err, nil)
		return
	}

	return
}

// AddUCASCourseID ...
func (m *Mongo) AddUCASCourseID(database, collection string, ucasCourseID *data.UCASCourseID) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Insert(ucasCourseID); err != nil {
		log.ErrorC("failed to create ucas course id resource", err, nil)
		return
	}

	return
}

// GetCAHCode ...
func (m *Mongo) GetCAHCode(database, collection, subjectCode string) (subjectObject *data.SubjectObject, err error) {
	s := m.Session.Copy()
	defer s.Close()

	if err = s.DB(database).C(collection).Find(bson.M{"code": subjectCode}).One(&subjectObject); err != nil {
		log.ErrorC("failed to find cah code resource", err, nil)
	}

	return
}

// DropCollection ...
func (m *Mongo) DropCollection(database, collection string) (err error) {
	s := m.Session.Copy()
	defer s.Close()

	if _, err = s.DB(database).C(collection).RemoveAll(nil); err != nil {
		return
	}

	return
}
