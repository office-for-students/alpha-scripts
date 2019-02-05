package main

import (
	"errors"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/handlers"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/mongo"
)

var (
	mongoURI string

	relativeFileLocation = "../files/"
	cahCodes             = "CAHCODES"
	commonData           = "COMMON"
	continuation         = "CONTINUATION"
	courseLocation       = "COURSELOCATION"
	degreeClass          = "DEGREECLASS"
	employment           = "EMPLOYMENT"
	entry                = "ENTRY"
	institution          = "INSTITUTION"
	institutionLocation  = "LOCATION"
	jobList              = "JOBLIST"
	jobType              = "JOBTYPE"
	leo                  = "LEO"
	nhsNSS               = "NHSNSS"
	nss                  = "NSS"
	qualifications       = "kisaims"
	salary               = "SALARY"
	subject              = "SBJ"
	tariff               = "TARIFF"
	ucasCourse           = "UCASCOURSEID"
)

var (
	wg sync.WaitGroup

	cahCodesCh            = make(chan int)
	courseLocationCh      = make(chan int)
	qualificationsCh      = make(chan int)
	commonDataCh          = make(chan int)
	continuationCh        = make(chan int)
	degreeClassCh         = make(chan int)
	employmentCh          = make(chan int)
	entryCh               = make(chan int)
	institutionCh         = make(chan int)
	institutionLocationCh = make(chan int)
	jobListCh             = make(chan int)
	jobTypeCh             = make(chan int)
	leoCh                 = make(chan int)
	nhsNSSCh              = make(chan int)
	nssCh                 = make(chan int)
	salaryCh              = make(chan int)
	subjectCh             = make(chan int)
	tariffCh              = make(chan int)
	ucasCourseCh          = make(chan int)
)

func main() {
	flag.StringVar(&mongoURI, "mongo-uri", mongoURI, "mongoDB URI")
	flag.StringVar(&relativeFileLocation, "relative-file-location", relativeFileLocation, "relative location of files")
	flag.Parse()

	if mongoURI == "" {
		log.Error(errors.New("missing mongo-uri flag"), nil)
		os.Exit(1)
	}

	mongodb := &mongo.Mongo{
		URI: mongoURI,
	}

	session, err := mongodb.Init()
	if err != nil {
		log.ErrorC("failed to initialise mongo", err, log.Data{"mongo_db": mongodb})
		os.Exit(1)
	}

	mongodb.Session = session

	common := handlers.Common{
		Mongo:                mongodb,
		RelativeFileLocation: relativeFileLocation,
	}

	go status()

	err = common.CreateCahCodes("courses", "cah-codes", cahCodes, cahCodesCh)
	if err != nil {
		log.ErrorC("Unsuccessfully attempted to load cah code data", err, nil)
		os.Exit(1)
	}

	wg.Add(18)

	go func() (err error) { // first goroutine as it has the largest dataset
		err = common.CreateJobList("statistics", "job-list", jobList, jobListCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateSubject("courses", "subjects", subject, subjectCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateCourseLocation("courses", "locations", courseLocation, courseLocationCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateQualifications("courses", "qualifications", qualifications, qualificationsCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateUCASCourseID("courses", "ucas-course-ids", ucasCourse, ucasCourseCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateInstitution("institutions", "raw", institution, institutionCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateInstitutionLocation("institutions", "locations", institutionLocation, institutionLocationCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateCommonData("statistics", "common", commonData, commonDataCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateContinuation("statistics", "continuation", continuation, continuationCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateDegreeClass("statistics", "degree-class", degreeClass, degreeClassCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateEmployment("statistics", "employment", employment, employmentCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateEntry("statistics", "entry", entry, entryCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateJobType("statistics", "job-type", jobType, jobTypeCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateLongitudinalEducationOutcomes("statistics", "leo", leo, leoCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateNHSNSS("statistics", "nhs-nss", nhsNSS, nhsNSSCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateNSS("statistics", "nss", nss, nssCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateSalary("statistics", "salary", salary, salaryCh)
		wg.Done()

		return
	}()

	go func() (err error) {
		err = common.CreateTariff("statistics", "tariff", tariff, tariffCh)
		wg.Done()

		return
	}()

	wg.Wait()

	// Allow for last log of data load before exiting script
	time.Sleep(1 * time.Second)
	if err != nil {
		log.ErrorC("Unsuccessfully attempted to load ofs data", err, nil)
		os.Exit(1)
	}

	log.Info("Successfully loaded ofs data", nil)
}

func status() {
	var (
		totalCount = 0

		cahCodeCount        = 0
		courseLocationCount = 0
		qualificationCount  = 0
		subjectCount        = 0
		ucasCourseIDCount   = 0

		institutionCount         = 0
		institutionLocationCount = 0

		commonDataCount   = 0
		continuationCount = 0
		degreeClassCount  = 0
		employmentCount   = 0
		entryCount        = 0
		jobListCount      = 0
		jobTypeCount      = 0
		leoCount          = 0
		nhsNSSCount       = 0
		nssCount          = 0
		salaryCount       = 0
		tariffCount       = 0
	)

	t := time.NewTicker(5 * time.Second)

	for {
		select {
		case n := <-cahCodesCh:
			cahCodeCount += n
			totalCount += n
		case n := <-courseLocationCh:
			courseLocationCount += n
			totalCount += n
		case n := <-qualificationsCh:
			qualificationCount += n
			totalCount += n
		case n := <-subjectCh:
			subjectCount += n
			totalCount += n
		case n := <-ucasCourseCh:
			ucasCourseIDCount += n
			totalCount += n
		case n := <-institutionCh:
			institutionCount += n
			totalCount += n
		case n := <-institutionLocationCh:
			institutionLocationCount += n
			totalCount += n
		case n := <-commonDataCh:
			commonDataCount += n
			totalCount += n
		case n := <-continuationCh:
			continuationCount += n
			totalCount += n
		case n := <-degreeClassCh:
			degreeClassCount += n
			totalCount += n
		case n := <-employmentCh:
			employmentCount += n
			totalCount += n
		case n := <-entryCh:
			entryCount += n
			totalCount += n
		case n := <-jobListCh:
			jobListCount += n
			totalCount += n
		case n := <-jobTypeCh:
			jobTypeCount += n
			totalCount += n
		case n := <-leoCh:
			leoCount += n
			totalCount += n
		case n := <-nhsNSSCh:
			nhsNSSCount += n
			totalCount += n
		case n := <-nssCh:
			nssCount += n
			totalCount += n
		case n := <-salaryCh:
			salaryCount += n
			totalCount += n
		case n := <-tariffCh:
			tariffCount += n
			totalCount += n
		case <-t.C:
			log.Info("Documents added",
				log.Data{
					"total":                 totalCount,
					"cah_codes":             cahCodeCount,
					"course_locations":      courseLocationCount,
					"qualifications":        qualificationCount,
					"subjects":              subjectCount,
					"ucas_course_ids":       ucasCourseIDCount,
					"institutions":          institutionCount,
					"institution_locations": institutionLocationCount,
					"common_datas":          commonDataCount,
					"continuations":         continuationCount,
					"degree_classes":        degreeClassCount,
					"employments":           employmentCount,
					"entries":               entryCount,
					"job_lists":             jobListCount,
					"job_types":             jobTypeCount,
					"leos":                  leoCount,
					"nhs_nsses":             nhsNSSCount,
					"nsses":                 nssCount,
					"salaries":              salaryCount,
					"tariffs":               tariffCount,
				},
			)
		}
	}
}
