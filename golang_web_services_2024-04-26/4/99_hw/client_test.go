package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

/*
1. Напишите в функции main код который просто
по фиксированным параметрам реализует логику SearchServer и выводит в консоль, без http
*/

func TestMyHttp(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	testCases := []struct {
		CaseName         string
		InputRequest     SearchRequest
		ExpectedResponse SearchResponse
	}{
		{
			CaseName: "find user by name",
			InputRequest: SearchRequest{
				Limit:      4,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResponse: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
		},
		{
			CaseName: "next page check",
			InputRequest: SearchRequest{
				Limit:  1,
				Offset: 0,
				Query:  "",
			},
			ExpectedResponse: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
		},
	}

	for _, tc := range testCases {
		resp, err := authSearchClient.FindUsers(tc.InputRequest)
		if err != nil {
			t.Errorf("[%s] unexpected error: %v", tc.CaseName, err)
			continue
		}
		if !reflect.DeepEqual(resp, &tc.ExpectedResponse) {
			t.Errorf("[%#v] wrong result, expected %#v, got %#v", tc.CaseName, tc.ExpectedResponse, resp)
		}
	}
}

func TestErrors(t *testing.T) {
	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         "",
	}
	resp, err := authSearchClient.FindUsers(SearchRequest{
		Limit: -1,
	})
	if resp != nil {
		t.Errorf("error: %v", resp)
	}
	if err == nil {
		t.Errorf("expected error, receive: %v", err)
	} else if err.Error() != "limit must be > 0" {
		t.Errorf("expected: %v, received: %v", "limit must be > 0", err.Error())
	}
}

func TestLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	resp, err := authSearchClient.FindUsers(SearchRequest{
		Limit: 26,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if len(resp.Users) != 25 {
		t.Errorf("expected 25, received: %v", len(resp.Users))
	}
	if resp.NextPage != true {
		t.Errorf("expected true, received: %v", resp.NextPage)
	}
}

func TestOffset(t *testing.T) {
	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         "",
	}

	resp, err := authSearchClient.FindUsers(SearchRequest{
		Offset: -1,
	})
	// сколько записей пропустить сначала
	if resp != nil {
		t.Errorf("expected nil, received: %v", resp)
	}
	if err == nil {
		t.Errorf("expected: error, received: %v", err)
	} else if err.Error() != "offset must be > 0" {
		t.Errorf("expected: %v, received: %v", "offset must be > 0", err)
	}

}

func TestUnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	ts.Close()

	resp, err := authSearchClient.FindUsers(SearchRequest{
		Limit:  1,
		Offset: 0,
	})
	if resp != nil {
		t.Fatalf("expected: nil, received: %v", resp)
	}
	if err == nil {
		t.Fatalf("expected: error, received: nil")
	}
	if !strings.HasPrefix(err.Error(), "unknown") {
		t.Errorf("expected: unknown, received: %v", err)
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	}))
	defer ts.Close()

	authSearchClient := SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}

	resp, err := authSearchClient.FindUsers(SearchRequest{})

	if resp != nil {
		t.Fatalf("expected: nil, received: %v", resp)
	}
	if err == nil {
		t.Fatal("expected: error, received: nil")
	}
	if !strings.HasPrefix(err.Error(), "timeout") {
		t.Errorf("expected: %v, received: %v", "timeout", err)
	}
}

func TestStatus(t *testing.T) {

	testCases := []struct {
		StatusCode    int
		ExpectedText  string
		CheckPrefix   bool
		ResponseError string
	}{
		{
			StatusCode:   401,
			ExpectedText: "Bad AccessToken",
		},
		{
			StatusCode:   500,
			ExpectedText: "SearchServer fatal error",
		},
		{
			StatusCode:   400,
			ExpectedText: "cant unpack error json:",
			CheckPrefix:  true,
		},
		{
			StatusCode:    400,
			ResponseError: "ErrorBadOrderField",
			ExpectedText:  "OrderField Height invalid",
		},
		{
			StatusCode:    400,
			ResponseError: "Unknown",
			ExpectedText:  "unknown bad request error: Unknown",
		},
		{
			StatusCode:    200, // как у нас покрылся этот тест case
			ResponseError: "",
			ExpectedText:  "cant unpack result json:",
			CheckPrefix:   true,
		},
	}

	for _, tc := range testCases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.StatusCode)
			if tc.ResponseError != "" { // проработать этот тест еще раз, так как плохо уложилось
				encoder := json.NewEncoder(w)
				err := encoder.Encode(SearchErrorResponse{Error: tc.ResponseError})
				if err != nil {
					t.Errorf("expected: nil, received: %v", err)
				}

			}
		}))
		defer ts.Close()

		authSearchClient := SearchClient{
			AccessToken: "",
			URL:         ts.URL,
		}
		resp, err := authSearchClient.FindUsers(SearchRequest{OrderField: "Height"})

		if resp != nil {
			t.Fatalf("expected: nil, received: %v", resp)
		}
		if err == nil {
			t.Fatalf("expected: error, received: %v", err)
		}
		if tc.CheckPrefix {
			if !strings.HasPrefix(err.Error(), tc.ExpectedText) {
				t.Errorf("expected: %v, received: %v", tc.ExpectedText, err)
			}
		} else if err.Error() != tc.ExpectedText {
			t.Fatalf("expected: %v, received: %v", tc.ExpectedText, err)
		}
	}
}
