package null

import (
	"database/sql/driver"
	"errors"
	"github.com/smgladkovskiy/structs"
	"github.com/smgladkovskiy/structs/zero"
	"time"
)

type Time struct {
	Time  time.Time
	Valid bool // isValid is true if Time is not NULL
}

// NewTime Создание Time переменной
func NewTime(v interface{}) (*Time, error) {
	var nt Time
	err := nt.Scan(v)
	return &nt, err
}

// Scan implements the Scanner interface for Time
func (nt *Time) Scan(value interface{}) error {
	switch v := value.(type) {
	case Time:
		*nt = v
		return nil
	case *Time:
		nt.Time, nt.Valid = v.Time, v.Valid
		return nil
	case zero.Time:
		tn := time.Time{}
		if v.Time == tn {
			return nil
		}
		nt.Time, nt.Valid = v.Time, true
		return nil
	case *zero.Time:
		tn := time.Time{}
		if v.Time == tn {
			return nil
		}
		nt.Time, nt.Valid = v.Time, true
		return nil
	case nil:
		*nt = Time{Time: time.Time{}, Valid: false}
		return nil
	case string:
		t, err := time.Parse(structs.TimeFormat(), v)
		if err != nil {
			*nt = Time{Time: time.Time{}, Valid: false}
			return err
		}
		*nt = Time{Time: t, Valid: true}
		return nil
	case time.Time:
		if v.IsZero() {
			*nt = Time{Time: time.Time{}, Valid: false}
			return nil
		}

		*nt = Time{Time: v, Valid: true}

		return nil
	case *time.Time:
		if v.IsZero() {
			*nt = Time{Time: time.Time{}, Valid: false}
			return nil
		}

		*nt = Time{Time: *v, Valid: true}

		return nil
	}

	return structs.TypeIsNotAcceptable{CheckedValue: value, CheckedType: nt}
}

// va implements the driver Valuer interface.
func (nt Time) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return structs.NullString, nil
	}
	if y := nt.Time.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(structs.TimeFormat())+2)
	b = append(b, '"')
	b = nt.Time.AppendFormat(b, structs.TimeFormat())
	b = append(b, '"')
	return b, nil
}

func (nt *Time) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return
	}
	nt.Time, err = time.Parse(structs.TimeFormat(), string(b))
	nt.Valid = err == nil
	return
}
