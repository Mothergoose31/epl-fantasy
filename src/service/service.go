package service

import (
	"encoding/json"
	"epl-fantasy/src/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type FPLService struct {
	BaseURL  string
	Endpoint string
}

func NewFPLService() (*FPLService, error) {
	err := godotenv.Load("URL.env")
	if err != nil {
		fmt.Printf("Warning: Error loading URL.env file: %v", err)
		fmt.Println("Falling back to system environment variables.")
	} else {
		fmt.Println("Successfully loaded URL.env file")
	}

	baseURL := os.Getenv("FPL_API_BASE_URL")
	endpoint := os.Getenv("FPL_BOOTSTRAP_ENDPOINT")

	return &FPLService{
		BaseURL:  baseURL,
		Endpoint: endpoint,
	}, nil
}

func getLatestWeek(events []config.Event) int {
	latestWeek := 0

	for _, event := range events {

		nameParts := strings.Fields(event.Name)
		if len(nameParts) != 2 {
			continue
		}

		weekStr := nameParts[1]
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			continue
		}

		if event.Finished && week > latestWeek {
			latestWeek = week
		}
	}

	return latestWeek
}

// ==================================================

func (s *FPLService) FetchFPLData() (*config.Data, []byte, error) {
	url := s.BaseURL + s.Endpoint

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %v", err)
	}

	var data config.Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, body, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")

	fmt.Println(data.GameWeek)

	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")

	fmt.Println(data.GameWeek)
	data.GameWeek = getLatestWeek(data.Events)

	return &data, body, nil
}
