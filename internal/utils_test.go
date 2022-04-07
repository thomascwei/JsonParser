package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var Body = []byte(`
{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      {
        "ff": 5878914884,
        "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460",
        "type": "thumbnail1",
        "inside": [
          {
            "aa": {"zz": 1.0},
            "bb": 2
          },
          {
            "aa": {"zz": 3.0},
            "bb": 4
          }
        ]
      },
      {
        "ff": 11,
        "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=4600",
        "type": "thumbnail2",
        "inside": [
          {
            "aa": {"zz": 11.11},
            "bb": 22
          },
          {
            "aa": {"zz": 33},
            "bb": 44
          }
        ]
      }
    ]
  },
  "company": {
    "name": "Acme",
    "list": [
      1,
      2.2,
      3.3
    ],
    "strs": [
      "a",
      "b",
      "c"
    ]
  }
}
`)

func TestCreateParseFunction(t *testing.T) {

	CreateParseFunction()
	t.Run("test case 1", func(t *testing.T) {
		if thisFunc, ok := EvalFuncHub["1"]; ok {
			got, err := thisFunc(Body)
			require.NoError(t, err)
			want := "Leonid"
			require.Equal(t, want, got[0])
			Trace.Println("test case 1 pass")
		}
	})
	t.Run("test case 2", func(t *testing.T) {
		if thisFunc, ok := EvalFuncHub["2"]; ok {
			got, err := thisFunc(Body)
			require.NoError(t, err)
			want := []string{"1",
				"2.2",
				"3.3"}
			require.Equal(t, want, got)
			Trace.Println("test case 2 pass")

		}
	})
	t.Run("test case 3", func(t *testing.T) {
		if thisFunc, ok := EvalFuncHub["3"]; ok {
			got, err := thisFunc(Body)
			require.NoError(t, err)
			want := "109"
			require.Equal(t, want, got[0])
			Trace.Println("test case 3 pass")

		}
	})
	t.Run("test case 4", func(t *testing.T) {
		if thisFunc, ok := EvalFuncHub["4"]; ok {
			got, err := thisFunc(Body)
			require.NoError(t, err)
			want := "109"
			require.Equal(t, want, got[0])
			Trace.Println("test case 4 pass")

		}
	})
	t.Run("test case 8", func(t *testing.T) {
		if thisFunc, ok := EvalFuncHub["8"]; ok {
			got, err := thisFunc(Body)
			require.NoError(t, err)
			want := []string{"a", "b", "c"}
			require.Equal(t, want, got)
			Trace.Println("test case 8 pass")

		}
	})
}

func TestGetOriAndObjectIDs(t *testing.T) {
	t.Run("test case 1", func(t *testing.T) {
		payload := EvalPayload{
			Body: Body,
			Id:   1,
		}
		ori, obj, err := GetOriAndObjectIDs(payload)
		require.NoError(t, err)
		require.Equal(t, "originObjectIdType0", ori)
		require.Equal(t, []string{"ObjectIdType0"}, obj)
	})
	t.Run("test case 2", func(t *testing.T) {
		payload := EvalPayload{
			Body: Body,
			Id:   2,
		}
		ori, obj, err := GetOriAndObjectIDs(payload)
		require.NoError(t, err)
		require.Equal(t, "originObjectIdType1", ori)
		require.Equal(t, []string{"aObjectIdType1-1-1", "aObjectIdType1-1-2", "aObjectIdType1-1-3"}, obj)
	})
	t.Run("test case 8", func(t *testing.T) {
		payload := EvalPayload{
			Body: Body,
			Id:   8,
		}
		ori, obj, err := GetOriAndObjectIDs(payload)
		require.NoError(t, err)
		require.Equal(t, "originObjectIdType0", ori)
		require.Equal(t, []string{"ObjectIdType1-2-1", "ObjectIdType1-2-2", "ObjectIdType1-2-3"}, obj)
	})
}

func TestInsertSensorDataPost(t *testing.T) {
	part := InsertSenorDataSingle{
		System:    "Test123",
		TimeStamp: int(time.Now().Unix()),
		TimeZone:  "",
		Data:      []InsideData{},
	}
	singleValue := InsideData{
		ObjectID:  "Test123",
		Value:     "000",
		TimeStamp: int(time.Now().Unix()),
	}
	part.Data = append(part.Data, singleValue)
	result := InsertSenorDataArray{part}
	resp, err := InsertSensorDataPost(result)
	require.NoError(t, err)
	Trace.Println(resp)
}
