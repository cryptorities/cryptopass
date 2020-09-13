package app

import (
	"encoding/base64"
	"time"
	"github.com/pkg/errors"
)

var ExecutableName = "cryptopass"

var ApplicationName = "CryptoPass"

var Copyright = "Copyright (C) Cryptorities LLC. All rights reserved."

var TimeDay = time.Hour * 24
var StartDate = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

var Encoding = base64.RawURLEncoding

var IssueSep = "*"
var RevokeSep = "/"

var DateFormatISO = "2006-01-02"


var ErrExpired = errors.New("EXPIRED")
