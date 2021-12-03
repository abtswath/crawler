package browser

import (
	"crawler/config"
	"crawler/utils"
	"github.com/go-rod/rod"
	"strconv"
	"strings"
	"time"
)

func (p *Page) fillForm() {
	p.wg.Add(1)
	go p.input()
	// TODO. select, textarea
}

func (p *Page) input() {
	defer p.wg.Done()
	elements, err := p.ElementsByJS(rod.Eval(`document.querySelector('input')`))
	if err != nil {
		return
	}
	for _, element := range elements {
		elementSetPossibleValue(element)
	}
}

func getPossibleValueByInputName(inputName string) string {
	if inputName == "" {
		return ""
	}
	for _, item := range config.PredictableInputValues {
		if utils.StringArrayInclude(item.Keyword, inputName) {
			return item.Value
		}
	}
	return ""
}

func elementSetPossibleValue(element *rod.Element) {
	inputType := getElementPropertyValue(element, "type", "text")
	inputName := getElementPropertyValue(element, "name", "")
	switch inputType {
	case InputTypeSearch:
		fallthrough
	case InputTypeText:
		possibleValue := getPossibleValueByInputName(inputName)
		if possibleValue == "" {
			possibleValue = getValidInputTextValue(element)
		}
		_ = element.Input(possibleValue)
	case InputTypeNumber:
		_ = element.Input(getValidInputNumberValue(element))
	case InputTypePassword:
		fallthrough
	case InputTypeEmail:
		fallthrough
	case InputTypeTel:
		fallthrough
	case InputTypeURL:
		_ = element.Input(config.PredictableInputValues[inputType].Value)
	case InputTypeDate:
		_ = element.Input(time.Now().Format("2006-01-02"))
	case InputTypeDatetimeLocal:
		_ = element.Input(time.Now().Format("2006-01-02 15:04:05"))
	case InputTypeMonth:
		_ = element.Input(time.Now().Format("2006-01"))
	case InputTypeTime:
		_ = element.Input(time.Now().Format("15:04"))
	case InputTypeWeek:
		_ = element.Input("1")
	// TODO. file, checkbox, radio
	}
}

func getElementPropertyValue(element *rod.Element, property string, defaultValue string) string {
	propertyValue, err := element.Property(property)
	if err != nil {
		return defaultValue
	}
	return strings.ToLower(propertyValue.String())
}

func getValidInputTextValue(element *rod.Element) string {
	length := 10
	var err error
	minLengthValue := getElementPropertyValue(element, "minlength", "")
	maxLengthValue := getElementPropertyValue(element, "maxlength", "")
	if minLengthValue != "" {
		length, err = strconv.Atoi(minLengthValue)
		if err != nil {
			length = 10
		}
	} else if maxLengthValue != "" {
		length, err = strconv.Atoi(maxLengthValue)
		if err != nil {
			length = 10
		}
	}
	return utils.RandomStr(length)
}

func getValidInputNumberValue(element *rod.Element) string {
	minValue := getElementPropertyValue(element, "min", "")
	maxValue := getElementPropertyValue(element, "max", "")
	if minValue != "" {
		return minValue
	}
	if maxValue != "" {
		return maxValue
	}
	return config.PredictableInputValues["number"].Value
}
