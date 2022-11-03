package target

import (
	"context"
	"errors"
	"fmt"
	"hotalert/logging"
	"hotalert/task"
	"io/ioutil"
	"net/http"
	"strings"
)

// ScrapeWebTask scraps the web page given the task.
func ScrapeWebTask(task *task.Task) error {
	// Parse options
	targetUrl, ok := task.Options["url"].(string)
	if !ok {
		logging.SugaredLogger.Errorf("Invalid task parameter url %v", targetUrl)
		return errors.New(fmt.Sprintf("Invalid task parameter url %v", targetUrl))
	}
	keywords, ok := task.Options["keywords"]
	if !ok {
		logging.SugaredLogger.Errorf("Invalid task parameter keywords %v", keywords)
		return errors.New(fmt.Sprintf("Invalid parameter keywords %v", keywords))
	}

	// Create a context with timeout specific to task.
	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()

	// Create a request with timeout.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)
	if err != nil {
		logging.SugaredLogger.Errorf("failed to build http request: %s", err)
		return err
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.SugaredLogger.Errorf("Failed to scrap page: %s", err)
		return err
	} else {
		if resp.StatusCode == 200 {
			pageBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logging.SugaredLogger.Errorf("Failed to read response from page. %s", err)
				return err
			}
			pageBodyStr := string(pageBody)

			// Search for matched keywords and save them.
			var matchedKeywords = make([]string, 0, 10)
			keywordsLst, ok := keywords.([]any)
			if !ok {
				logging.SugaredLogger.Errorf("Invalid task parameter keywords %v", keywords)
				return errors.New(fmt.Sprintf("Invalid parameter keywords %v", keywords))
			}
			for _, value := range keywordsLst {
				valueStr, ok := value.(string)
				if !ok {
					logging.SugaredLogger.Errorf("Invalid value in task keywords, not a string %v", valueStr)
					return errors.New(fmt.Sprintf("Invalid value in task keywords, not a string %v", valueStr))
				}
				if strings.Contains(pageBodyStr, valueStr) {
					matchedKeywords = append(matchedKeywords, valueStr)
				}
			}

			// If we have matched keywords post an alert.
			if len(matchedKeywords) > 0 {
				task.Alerter.PostAlert(context.Background(), matchedKeywords)
			}
		} else {
			logging.SugaredLogger.Errorf("Failed to query website, status code %d", resp.StatusCode)
			return errors.New(fmt.Sprintf("Failed to query website, status code %d", resp.StatusCode))
		}
	}
	return nil
}
