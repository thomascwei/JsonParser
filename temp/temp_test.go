package thomas

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var data = []byte(`{
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
`    )
func TestFunc1(t *testing.T) {

	rr, err := func1(data)
	require.NoError(t, err)
	fmt.Println(rr)
}

func TestFunc2(t *testing.T) {

	rr, err := func2(data)
	require.NoError(t, err)
	fmt.Println(rr)
}

func TestFunc8(t *testing.T) {

	rr, err := func8(data)
	require.NoError(t, err)
	fmt.Println(rr)
}

func TestFunc9(t *testing.T) {

	rr, err := func9(data)
	require.NoError(t, err)
	fmt.Println(rr)
}