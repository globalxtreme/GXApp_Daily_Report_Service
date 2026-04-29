package core

import (
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm/schema"
	error2 "service/internal/pkg/error"
	"strconv"
	"strings"
	"time"
)

type NumberOptions struct {
	FullNumber bool
	PadLength  int
	RandMin    int
	RandMax    int
}

type NumberGeneratorInterface interface {
	Generate() string
	Prefix() string
}

func ModelNumber(mdl schema.Tabler, generator NumberGeneratorInterface, opt ...NumberOptions) string {
	padLength := 4
	fullNumber := true
	randMin := 111
	randMax := 999

	if len(opt) > 0 {
		padLength = opt[0].PadLength
		fullNumber = opt[0].FullNumber
		randMin = opt[0].RandMin
		randMax = opt[0].RandMax
	}

	increment := GetIncrementMonthly(mdl)

	number := time.Now().Format("010602")
	number += StrPadLeft(strconv.Itoa(int(increment)), padLength, '0')
	number += strconv.Itoa(RandInt(randMin, randMax))

	if !fullNumber {
		return number
	}

	return strings.ToUpper(generator.Prefix() + number)
}

func RedisNumber(prefix string, opt ...NumberOptions) string {
	padLength := 7
	fullNumber := true

	if len(opt) > 0 {
		padLength = opt[0].PadLength
		fullNumber = opt[0].FullNumber
	}

	if !IsProduction() {
		prefix = "DEV" + prefix
	}

	date := time.Now().Format("20060102")
	key := fmt.Sprintf("number_counter:%s:%s", prefix, date)

	conn := xtremepkg.RedisPool.Get()
	defer conn.Close()

	seq, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		error2.ErrXtremeINCRNumber(err.Error())
	}

	// Exp: 24 Hours
	conn.Do("EXPIRE", key, 86400)

	format := fmt.Sprintf("%s%%0%dd", date, padLength)
	if !fullNumber {
		return fmt.Sprintf(format, seq)
	}

	return fmt.Sprintf("%s"+format, prefix, seq)
}
