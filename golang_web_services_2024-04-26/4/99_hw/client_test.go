package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
