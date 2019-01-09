package statistics

import (
	"sync"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
)

type statConfig struct {
	kisCourseID string
	kisMode     string
	publicUKPRN string
	uri         string
}

// Get ...
func Get(mongoURI, publicUKPRN, kisCourseID, kisMode string) (*data.Statistics, error) {
	stat := statConfig{
		kisCourseID: kisCourseID,
		kisMode:     kisMode,
		publicUKPRN: publicUKPRN,
		uri:         mongoURI,
	}

	var wg sync.WaitGroup
	var (
		continuation []*data.Continuation
		employment   []*data.Employment
		jobList      data.JobList
		jobType      []*data.JobType
		leo          []*data.LEO
		salary       []*data.Salary
	)

	wg.Add(6)
	go func() {
		continuation, _ = stat.continuation()
		wg.Done()

		return
	}()

	go func() {
		employment, _ = stat.employment()
		wg.Done()

		return
	}()

	go func() {
		jobList, _ = stat.jobList()
		wg.Done()

		return
	}()

	go func() {
		jobType, _ = stat.jobType()
		wg.Done()

		return
	}()

	go func() {
		leo, _ = stat.leo()
		wg.Done()

		return
	}()

	go func() {
		salary, _ = stat.salary()
		wg.Done()

		return
	}()

	wg.Wait()

	stats := &data.Statistics{
		Continuation: continuation,
		Employment:   employment,
		JobList:      &jobList,
		JobType:      jobType,
		LEO:          leo,
		Salary:       salary,
	}

	return stats, nil
}

func (stat *statConfig) continuation() (results []*data.Continuation, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("statistics").C("continuation").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find continuation resources for course", err, nil)
	}

	return
}

func (stat *statConfig) employment() (results []*data.Employment, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("statistics").C("employment").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find employment resources for course", err, nil)
	}

	return
}

func (stat *statConfig) jobList() (results data.JobList, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var jobs []data.JobOrder

	if err = session.DB("statistics").C("job-list").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&jobs); err != nil {
		log.ErrorC("failed to find job list resources for course", err, nil)
		return
	}

	m := make(map[int][]data.Job)
	for _, job := range jobs {
		newJob := data.Job{
			Job:                  job.Job,
			PercentageOfStudents: job.PercentageOfStudents,
		}
		m[job.Order] = append(m[job.Order], newJob)
	}

	var items []data.JobItem
	for i := 1; i <= len(m); i++ {
		item := data.JobItem{
			List:  m[i],
			Order: i,
		}
		items = append(items, item)
	}
	results.Items = items

	var common data.Common

	// Get metadata for stats (common), e.g. aggregation level, response rate and number of students
	if err = session.DB("statistics").C("common").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).One(&common); err != nil {
		log.ErrorC("failed to find job list resources for course", err, nil)
		return
	}

	results.AggregationLevel = common.AggregationLevel
	results.NumberOfStudents = common.NumberOfStudents
	results.ResponseRate = common.ResponseRate
	results.SubjectCode = common.SubjectCode
	results.Unavailable = common.Unavailable

	return
}

func (stat *statConfig) jobType() (jobTypes []*data.JobType, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("statistics").C("job-type").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&jobTypes); err != nil {
		log.ErrorC("failed to find job type resources for course", err, nil)
	}

	return
}

func (stat *statConfig) leo() (results []*data.LEO, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("statistics").C("leo").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find leo resources for course", err, nil)
	}

	return
}

func (stat *statConfig) salary() (salary []*data.Salary, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var results []*data.RawSalary
	if err = session.DB("statistics").C("salary").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find salary resources for course", err, nil)
	}

	for _, result := range results {
		s := &data.Salary{
			AggregationLevel:  result.AggregationLevel,
			NumberOfGraduates: result.NumberOfStudents,
			ResponseRate:      result.ResponseRate,
			SubjectCode:       result.SubjectCode,
			Unavailable:       result.Unavailable,
		}

		if result.InstitutionCourseSalarySixMonthsAfterGraduation != nil {
			s.LowerQuartileRange = result.InstitutionCourseSalarySixMonthsAfterGraduation.LowerQuartile
			s.Median = result.InstitutionCourseSalarySixMonthsAfterGraduation.Median
			s.HigherQuartileRange = result.InstitutionCourseSalarySixMonthsAfterGraduation.UpperQuartile
		}

		salary = append(salary, s)
	}

	return
}
