package statistics

import (
	"strconv"
	"sync"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
)

type statConfig struct {
	countryCode string
	kisCourseID string
	kisMode     string
	publicUKPRN string
	uri         string
}

var reason = map[int]string{
	0: "There is not enough data available to publish for this course. This is either because the course is small or we have not had enough survey responses. **This does not reflect on the quality of the course.**",
	1: "There is no data available for this course, as the course has either not run yet, or has not been running long enough for this data to be available.  **This does not reflect on the quality of the course.**",
	2: "There is no data available for this course. **This does not reflect on the quality of the course.**",
	3: "There was not enough data to publish information specifically for this course. This is either because the course size is small or not enough students responded to a survey. For this reason, the data displayed is for all students in ",
	4: "There is no data available for this course. This is because the course has not yet run or has not been running long enough for this data to be available. For this reason, the data displayed is for students on other courses in ",
	5: "Data for students in the last two years of this course has been combined, as there was not enough data to publish information for last year only.",
	6: "There is no data available for the subject area of this course. This may be because we only have data for a small number of students or because we do not yet have data. **This does not reflect on the quality of the course.**",
	7: "We only have this data for English universities and colleges. This is because of differences in either policy or legislation relating to this data in the other countries of the UK. **This does not reflect on the quality of the course.**",
}

// Get ...
func Get(mongoURI, publicUKPRN, kisCourseID, kisMode, countryCode string) (*data.Statistics, *data.Subject, error) {
	stat := statConfig{
		countryCode: countryCode,
		kisCourseID: kisCourseID,
		kisMode:     kisMode,
		publicUKPRN: publicUKPRN,
		uri:         mongoURI,
	}

	var wg sync.WaitGroup
	var (
		continuation []*data.Continuation
		employment   []*data.Employment
		jobList      *data.JobList
		jobType      []*data.JobType
		leo          []*data.LEO
		salary       []*data.Salary
		subject      *data.Subject
	)

	wg.Add(7)
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

	go func() {
		subject, _ = stat.subject()
		wg.Done()

		return
	}()

	wg.Wait()

	stats := &data.Statistics{
		Continuation: continuation,
		Employment:   employment,
		JobList:      jobList,
		JobType:      jobType,
		LEO:          leo,
		Salary:       salary,
	}

	return stats, subject, nil
}

func (stat *statConfig) subject() (subject *data.Subject, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var subjectObject *data.SubjectItem
	if err = session.DB("courses").C("subjects").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).One(&subjectObject); err != nil {
		log.ErrorC("failed to find subject resource for course", err, nil)
	}

	subject = &data.Subject{
		Code: subjectObject.Subject.Code,
		Name: subjectObject.Subject.Name,
	}

	return
}

func (stat *statConfig) continuation() (continuations []*data.Continuation, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var results []*data.ContinuationRaw
	if err = session.DB("statistics").C("continuation").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find continuation resources for course", err, nil)
	}

	for _, result := range results {
		continuation := &data.Continuation{
			AggregationLevel:             result.AggregationLevel,
			ContinuingWithProvider:       result.ContinuingWithProvider,
			Dormant:                      result.Dormant,
			GainingIntendedAwardOrHigher: result.GainingIntendedAwardOrHigher,
			GainedLowerAward:             result.GainedLowerAward,
			LeavingCourse:                result.LeavingCourse,
			Subject:                      result.Subject,
		}

		subjectName := ""
		if result.Subject != nil {
			subjectName = result.Subject.Name
		}

		if result.AggregationLevel != 0 {
			continuation.Unavailable = handleDelhiUnavailableEnum(true, result.AggregationLevel, result.Unavailable, subjectName)
		} else {
			continuation.Unavailable = handleDelhiUnavailableEnum(false, result.AggregationLevel, result.Unavailable, subjectName)
		}

		continuations = append(continuations, continuation)
	}

	return
}

func (stat *statConfig) employment() (employments []*data.Employment, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var results []*data.EmploymentRaw
	if err = session.DB("statistics").C("employment").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find employment resources for course", err, nil)
	}

	for _, result := range results {
		employment := &data.Employment{
			AggregationLevel:           result.AggregationLevel,
			AssumedToBeUnemployed:      result.AssumedToBeUnemployed,
			NumberOfStudents:           result.NumberOfStudents,
			InStudy:                    result.InStudy,
			InWork:                     result.InWork,
			InWorkAndStudy:             result.InWorkAndStudy,
			InWorkOrStudy:              result.InWorkOrStudy,
			NotAvailableForWorkOrStudy: result.NotAvailableForWorkOrStudy,
			ResponseRate:               result.ResponseRate,
			Subject:                    result.Subject,
		}

		subjectName := ""
		if result.Subject != nil {
			subjectName = result.Subject.Name
		}

		if result.AggregationLevel != 0 {
			employment.Unavailable = handleDelhiUnavailableEnum(true, result.AggregationLevel, result.Unavailable, subjectName)
		} else {
			employment.Unavailable = handleDelhiUnavailableEnum(false, result.AggregationLevel, result.Unavailable, subjectName)
		}

		employments = append(employments, employment)
	}

	return
}

func (stat *statConfig) jobList() (*data.JobList, error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return nil, err
	}
	defer session.Close()

	var jobs []data.JobOrder

	if err = session.DB("statistics").C("job-list").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&jobs); err != nil {
		log.ErrorC("failed to find job list resources for course", err, nil)
		return nil, err
	}

	m := make(map[string][]data.Job)
	for _, job := range jobs {
		newJob := data.Job{
			Job:                  job.Job,
			Order:                job.Order,
			PercentageOfStudents: job.PercentageOfStudents,
		}

		var subjectCode string
		if job.Subject != nil {
			newJob.Subject = nil
			subjectCode = job.Subject.Code
		}

		m[subjectCode] = append(m[subjectCode], newJob)
	}

	results := &data.JobList{}
	for subject, jobs := range m {
		item := &data.SubjectItem{
			List: jobs,
		}

		var selector bson.M
		if subject == "" {
			selector = bson.M{
				"public_ukprn":  stat.publicUKPRN,
				"kis_course_id": stat.kisCourseID,
				"kis_mode":      stat.kisMode,
			}
		} else {
			selector = bson.M{
				"public_ukprn":  stat.publicUKPRN,
				"kis_course_id": stat.kisCourseID,
				"kis_mode":      stat.kisMode,
				"subject.code":  subject,
			}
		}

		var common data.Common

		// Get metadata for stats (common), e.g. aggregation level, response rate and number of students
		if err = session.DB("statistics").C("common").Find(selector).One(&common); err != nil {
			log.ErrorC("failed to find job list resources for course", err, nil)
			return nil, err
		}

		item.AggregationLevel = common.AggregationLevel
		item.NumberOfStudents = common.NumberOfStudents
		item.ResponseRate = common.ResponseRate
		item.Subject = common.Subject

		if common.Subject != nil {
			item.Unavailable = handleDelhiUnavailableEnum(true, common.AggregationLevel, common.Unavailable, common.Subject.Name)
		}

		results.Items = append(results.Items, item)
	}

	if len(results.Items) < 1 {
		var common data.Common

		// Get metadata for stats (common), e.g. aggregation level, response rate and number of students
		if err = session.DB("statistics").C("common").Find(bson.M{
			"public_ukprn":  stat.publicUKPRN,
			"kis_course_id": stat.kisCourseID,
			"kis_mode":      stat.kisMode,
		}).One(&common); err != nil {
			log.ErrorC("failed to find job list resources for course", err, nil)
			return nil, err
		}

		results.Unavailable = handleDelhiUnavailableEnum(false, common.AggregationLevel, common.Unavailable, "")
	}

	return results, nil
}

func (stat *statConfig) jobType() (jobTypes []*data.JobType, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var results []*data.JobTypeRaw
	if err = session.DB("statistics").C("job-type").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find job type resources for course", err, nil)
	}

	for _, result := range results {
		jobType := &data.JobType{
			AggregationLevel:                result.AggregationLevel,
			NonProfessionalOrManagerialJobs: result.NonProfessionalOrManagerialJobs,
			NumberOfStudents:                result.NumberOfStudents,
			ProfessionalOrManagerialJobs:    result.ProfessionalOrManagerialJobs,
			ResponseRate:                    result.ResponseRate,
			Subject:                         result.Subject,
			UnknownProfessions:              result.UnknownProfessions,
		}

		subjectName := ""
		if result.Subject != nil {
			subjectName = result.Subject.Name
		}

		if result.AggregationLevel != 0 {
			jobType.Unavailable = handleDelhiUnavailableEnum(true, result.AggregationLevel, result.Unavailable, subjectName)
		} else {
			jobType.Unavailable = handleDelhiUnavailableEnum(false, result.AggregationLevel, result.Unavailable, subjectName)
		}

		jobTypes = append(jobTypes, jobType)
	}

	return
}

func (stat *statConfig) leo() (leos []*data.LEO, err error) {
	session, err := mgo.Dial(stat.uri)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	var results []*data.LEORaw
	if err = session.DB("statistics").C("leo").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find leo resources for course", err, nil)
	}

	for _, result := range results {
		leo := &data.LEO{
			AggregationLevel:    result.AggregationLevel,
			HigherQuartileRange: result.HigherQuartileRange,
			LowerQuartileRange:  result.LowerQuartileRange,
			Median:              result.Median,
			NumberOfGraduates:   result.NumberOfGraduates,
			Subject:             result.Subject,
		}

		if leo.HigherQuartileRange == 0 {
			leo.Unavailable = handleLEOUnavailableEnum(stat.countryCode, result.Unavailable)
		}

		leos = append(leos, leo)
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

	var results []*data.SalaryRaw
	if err = session.DB("statistics").C("salary").Find(bson.M{"public_ukprn": stat.publicUKPRN, "kis_course_id": stat.kisCourseID, "kis_mode": stat.kisMode}).All(&results); err != nil {
		log.ErrorC("failed to find salary resources for course", err, nil)
	}

	for _, result := range results {
		s := &data.Salary{
			AggregationLevel:  result.AggregationLevel,
			NumberOfGraduates: result.NumberOfStudents,
			ResponseRate:      result.ResponseRate,
			Subject:           result.Subject,
		}

		subjectName := ""
		if result.Subject != nil {
			subjectName = result.Subject.Name
		}

		if result.SubjectSalarySixMonthsAfterGraduation != nil && result.SubjectSalarySixMonthsAfterGraduation.LowerQuartile != 0 {
			s.LowerQuartileRange = result.SubjectSalarySixMonthsAfterGraduation.LowerQuartile
			s.Median = result.SubjectSalarySixMonthsAfterGraduation.Median
			s.HigherQuartileRange = result.SubjectSalarySixMonthsAfterGraduation.UpperQuartile

			s.Unavailable = handleDelhiUnavailableEnum(true, result.AggregationLevel, result.Unavailable, subjectName)
		} else {
			s.Unavailable = handleDelhiUnavailableEnum(false, result.AggregationLevel, result.Unavailable, subjectName)
		}

		salary = append(salary, s)
	}

	return
}

func handleDelhiUnavailableEnum(hasData bool, aggregationLevel int, unavailable, subjectName string) *data.Unavailable {

	if aggregationLevel == 14 {
		return nil
	}

	unavailableObject := &data.Unavailable{}
	if hasData {
		switch unavailable {
		case "0":
			if aggregationLevel < 20 {
				unavailableObject.Reason = reason[3] + subjectName + "."
			} else {
				unavailableObject.Reason = reason[3] + subjectName + " across the last two years."
			}
		case "1":
			if aggregationLevel < 20 {
				unavailableObject.Reason = reason[4] + subjectName + "."
			} else {
				if aggregationLevel == 24 {
					unavailableObject.Reason = reason[5]
				} else {
					unavailableObject.Reason = reason[4] + subjectName + " across the last two years."
				}
			}
		}
	}

	if !hasData {
		switch unavailable {
		case "0":
			unavailableObject.Reason = reason[0]
		case "1":
			unavailableObject.Reason = reason[1]
		case "2":
			unavailableObject.Reason = reason[2]
		}
	}

	var err error
	unavailableObject.Code, err = strconv.Atoi(unavailable)
	if err != nil {
		log.ErrorC("unavailable code invalid", err, log.Data{"code": unavailable})
		unavailableObject.Code = 3
	}

	return unavailableObject
}

func handleLEOUnavailableEnum(countryCode, unavailable string) *data.Unavailable {

	unavailableObject := &data.Unavailable{}
	switch countryCode {
	case "XF": // England
		unavailableObject.Reason = reason[6]
	case "XG": // Northern Ireland
		unavailableObject.Reason = reason[7]
	case "XH": // Scotland
		unavailableObject.Reason = reason[7]
	case "XI": // Wales
		unavailableObject.Reason = reason[7]
	}

	var err error
	unavailableObject.Code, err = strconv.Atoi(unavailable)
	if err != nil {
		log.ErrorC("unavailable code invalid", err, log.Data{"code": unavailable})
		unavailableObject.Code = 3
	}

	return unavailableObject
}
