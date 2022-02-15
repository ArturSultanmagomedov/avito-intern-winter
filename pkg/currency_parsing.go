package pkg

import (
	"encoding/json"
	"for_avito_tech_with_gin/pkg/model"
	"for_avito_tech_with_gin/pkg/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

//go:generate mockgen -source=currency_parsing.go -destination=mocks/mock.go

type CurrencyCalculator interface {
	UpdateRates()
	ConvertRubTo(currency string, sum float32) (float64, error)
}

var list model.CurrencyList

type DefaultCurrencyCalculator struct{}

func (r *DefaultCurrencyCalculator) UpdateRates() {
	updateCurrencyJson()
}

func (r *DefaultCurrencyCalculator) ConvertRubTo(currency string, sum float32) (float64, error) {
	logrus.Debugf("ConvertRubTo invoke, currency = %s, sum = %f", currency, sum)
	nominal, ok1 := list.List[currency]["Nominal"]
	value, ok2 := list.List[currency]["Value"]
	if !ok1 || !ok2 {
		return 0, &service.WrongParam{Param: "currency"}
	}

	v, err := getFloat(value)
	if err != nil {
		logrus.Error(err)
		return 0, &service.InternalServerError{}
	}
	n, err := getFloat(nominal)
	if err != nil {
		logrus.Error(err)
		return 0, &service.InternalServerError{}
	}
	logrus.Debugf("nominal = %f, value = %f", nominal, value)

	return float64(sum) / v * n, nil
}

func updateCurrencyJson() {
	logrus.Debugf("updateCurrencyJson invoke")
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		logrus.Error(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		logrus.Error(err)
	}
}

func getFloat(unk interface{}) (float64, error) {
	floatType := reflect.TypeOf(float64(0))
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, errors.New("412, wrong currency param")
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
