package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
)

var (
	mongoURI string
	wg       sync.WaitGroup

	readCountCh    = make(chan int)
	failedCourseCh = make(chan int)
	failedURLCh    = make(chan int)
	sem            = make(chan int, 5)

	mongoDatabase   = "courses"
	mongoCollection = "courses"
	mongoSize       = 500

	filename = "broken-urls.csv"
)

func main() {
	flag.StringVar(&mongoURI, "mongo-uri", mongoURI, "mongoDB URI")
	flag.Parse()

	if mongoURI == "" {
		log.Error(errors.New("missing mongo-uri flag"), nil)
		os.Exit(1)
	}

	s, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("error creating mongo db session", err, log.Data{"filename": filename})
	}

	if err := os.Remove(filename); err != nil {
		if err.Error() != "remove "+filename+": no such file or directory" {
			log.ErrorC("error removing file", err, log.Data{"filename": filename, "errorMessage": err.Error()})
			os.Exit(1)
		}
	}

	connection, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.ErrorC("error opening file", err, log.Data{"filename": filename})
	}

	writeToFile(connection, filename, "ukprn,kis_mode,kis_course_id,url_field,url_value,status_code,status_description")

	go status()

	it := s.DB(mongoDatabase).C(mongoCollection).Find(bson.M{}).Batch(mongoSize).Iter()

	for {
		courses := make([]data.Course, mongoSize)

		itx := 0
		for ; itx < len(courses); itx++ {
			result := data.Course{}

			if !it.Next(&result) {
				break
			}
			courses[itx] = result
		}

		if itx == 0 {
			// time.Sleep(10 * time.Second)
			break
		}

		// This will block if we've reached our concurrecy limit (sem buffer size)
		iterateCourses(connection, courses, itx)
	}

	// May need to sleep here
	wg.Wait()

	if err := connection.Close(); err != nil {
		log.ErrorC("error closing file: "+filename, err, nil)
	}
}

func iterateCourses(c *os.File, courses []data.Course, length int) {
	// Wait on semaphore if we've reached our concurrency limit
	wg.Add(1)
	sem <- 1

	go func() {
		defer func() {
			<-sem
			wg.Done()
		}()

		i := 0
		// loop through course documents
		for i < length {
			course := courses[i]

			checkCourseURLs(c, course)

			readCountCh <- 1
			i++
		}
	}()
}

func checkCourseURLs(c *os.File, course data.Course) {

	// "ukprn,kis_mode,kis_course_id,url_field,url_value,status_code,status_description"
	var linkFailure bool
	if course.Links != nil {
		if course.Links.AssessmentMethod != nil {
			if course.Links.AssessmentMethod.English != "" {
				linkName := "ASSURL"
				if failed := makeRequest(c, course, linkName, course.Links.AssessmentMethod.English); failed {
					linkFailure = true
				}
			}

			if course.Links.AssessmentMethod.Welsh != "" {
				linkName := "ASSURLW"
				if failed := makeRequest(c, course, linkName, course.Links.AssessmentMethod.Welsh); failed {
					linkFailure = true
				}
			}
		}

		if course.Links.CoursePage != nil {
			if course.Links.CoursePage.English != "" {
				linkName := "CRSEURL"
				if failed := makeRequest(c, course, linkName, course.Links.CoursePage.English); failed {
					linkFailure = true
				}
			}

			if course.Links.CoursePage.Welsh != "" {
				linkName := "CRSEURLW"
				if failed := makeRequest(c, course, linkName, course.Links.CoursePage.Welsh); failed {
					linkFailure = true
				}
			}
		}

		if course.Links.EmploymentDetails != nil {
			if course.Links.EmploymentDetails.English != "" {
				linkName := "EMPLOYURL"
				if failed := makeRequest(c, course, linkName, course.Links.EmploymentDetails.English); failed {
					linkFailure = true
				}
			}

			if course.Links.EmploymentDetails.Welsh != "" {
				linkName := "EMPLOYURLW"
				if failed := makeRequest(c, course, linkName, course.Links.EmploymentDetails.Welsh); failed {
					linkFailure = true
				}
			}
		}

		if course.Links.FinancialSupport != nil {
			if course.Links.FinancialSupport.English != "" {
				linkName := "SUPPORTURL"
				if failed := makeRequest(c, course, linkName, course.Links.FinancialSupport.English); failed {
					linkFailure = true
				}
			}

			if course.Links.FinancialSupport.Welsh != "" {
				linkName := "SUPPORTURLW"
				if failed := makeRequest(c, course, linkName, course.Links.FinancialSupport.Welsh); failed {
					linkFailure = true
				}
			}
		}

		if course.Links.LearningAndTeaching != nil {
			if course.Links.LearningAndTeaching.English != "" {
				linkName := "LTURL"
				if failed := makeRequest(c, course, linkName, course.Links.LearningAndTeaching.English); failed {
					linkFailure = true
				}
			}

			if course.Links.LearningAndTeaching.Welsh != "" {
				linkName := "LTURLW"
				if failed := makeRequest(c, course, linkName, course.Links.LearningAndTeaching.Welsh); failed {
					linkFailure = true
				}
			}
		}
	}

	if linkFailure {
		failedCourseCh <- 1
	}
}

func makeRequest(c *os.File, course data.Course, linkName, path string) (linkFailure bool) {
	logData := log.Data{"path": path}

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorC("failed to parse url", err, logData)
	}
	path = URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.ErrorC("failed to create request", err, logData)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
		// Timeout: time.Duration(30 * time.Second),
	}

	r, err := client.Do(req)
	if err != nil {
		log.ErrorC("failed to make request", err, logData)

		failedURLCh <- 1
		linkFailure = true

		// Assume i've been redirected
		line := course.Institution.UKPRN + "," + course.KISCourseID + "," + course.Mode.Code + "," + linkName + "," + path + "," + "500,internal server error"
		writeToFile(c, filename, line)
		return
	}
	defer r.Body.Close()

	if r.StatusCode == 301 || r.StatusCode == 302 {
		ok := checkHTTPS(path)
		// log.Debug("check ok value", log.Data{"ok": ok})
		if !ok {
			failedURLCh <- 1
			linkFailure = true

			statusDescription := r.Status
			line := course.Institution.UKPRN + "," + course.KISCourseID + "," + course.Mode.Code + "," + linkName + "," + path + "," + strconv.Itoa(r.StatusCode) + "," + statusDescription
			writeToFile(c, filename, line)
		}
	}

	if r.StatusCode > 302 {

		failedURLCh <- 1
		linkFailure = true

		statusDescription := r.Status
		line := course.Institution.UKPRN + "," + course.KISCourseID + "," + course.Mode.Code + "," + linkName + "," + path + "," + strconv.Itoa(r.StatusCode) + "," + statusDescription
		writeToFile(c, filename, line)
	}

	return
}

var captureHTTP = regexp.MustCompile(`^(http)(://.*)$`)

func checkHTTPS(path string) (success bool) {
	pathComponents := captureHTTP.FindStringSubmatch(path)
	logData := log.Data{"func": "checkHTTPS", "path_components": pathComponents}
	// log.Debug("path components?", log.Data{"path_components": pathComponents})
	if len(pathComponents) < 3 {
		log.Debug("is not http", logData)
		return
	}

	if pathComponents[1] == "https" {
		log.Debug("is already https", logData)
		return
	}

	httpsPath := "https" + pathComponents[2]
	logData["new_path"] = httpsPath

	URL, err := url.Parse(httpsPath)
	if err != nil {
		log.ErrorC("failed to parse url", err, logData)
		return
	}
	path = URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.ErrorC("failed to create request", err, logData)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
		// Timeout: time.Duration(30 * time.Second),
	}

	r, err := client.Do(req)
	if err != nil {
		log.ErrorC("failed to make request", err, logData)
		return
	}
	defer r.Body.Close()

	if r.StatusCode < 300 {
		success = true
	}

	return
}

func status() {
	var (
		readCourseCounter   = 0
		failedCourseCounter = 0
		failedURLCounter    = 0
	)

	t := time.NewTicker(time.Second)

	for {
		select {
		case n := <-readCountCh:
			readCourseCounter += n
		case n := <-failedCourseCh:
			failedCourseCounter += n
		case n := <-failedURLCh:
			failedURLCounter += n
		case <-t.C:
			line := fmt.Sprintf("Courses read: %v  Course-Failure: %v  URL-Failure: %v", readCourseCounter, failedCourseCounter, failedURLCounter)
			log.Info(line, nil)
		}
	}
}

func writeToFile(connection *os.File, location string, line string) {
	_, err := connection.WriteString(line + "\n")
	if err != nil {
		log.ErrorC("error writing to file", err, log.Data{"file_location": location, "line": line})
	}
}
