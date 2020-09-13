package util_test

import (
	"github.com/cryptorities/cryptopass/pkg/app"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDays(t *testing.T) {

	tm := util.ParseDaysOffset(0)
	s := tm.Format(app.DateFormatISO)
	assert.Equal(t, "2020-01-01", s)
	assert.Equal(t, 0, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(1)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2020-01-02", s)
	assert.Equal(t, 1, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(-1)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2019-12-31", s)
	assert.Equal(t, -1, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(100)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2020-04-10", s)
	assert.Equal(t, 100, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(-100)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2019-09-23", s)
	assert.Equal(t, -100, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(1000)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2022-09-27", s)
	assert.Equal(t, 1000, util.DaysOffset(&tm))

	tm = util.ParseDaysOffset(10000)
	s = tm.Format(app.DateFormatISO)
	assert.Equal(t, "2047-05-19", s)
	assert.Equal(t, 10000, util.DaysOffset(&tm))

}
