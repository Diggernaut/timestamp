package timestamp

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {

	if y := time.Time(t).Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(time.Time(t).Format(`"` + time.RFC3339Nano + `"`)), nil

}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var tim time.Time
	var err error

	if string(b) != "\"\"" { // check empty string
		tim, err = time.Parse(`"`+time.RFC3339+`"`, string(b))
		*t = Timestamp(tim)
		if err != nil {
			*t = Timestamp(tim) // return 0 time
		}
	}

	return nil
}

func (t Timestamp) GetBSON() (interface{}, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}

	return time.Time(t), nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var tm time.Time

	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}

	*t = Timestamp(tm)

	return nil
}

func (t *Timestamp) IsZeroTimestamp() bool {
	return t.GetTime().IsZero()
}

func (t *Timestamp) Format(layout string) string {
	return time.Time(*t).Format(layout)
}
func (t Timestamp) String() string {
	return t.Format(time.RFC3339)
}

func (t *Timestamp) GetTime() time.Time {
	return time.Time(*t)
}

func GetTimeStamp(t time.Time) Timestamp {
	return Timestamp(t)
}
