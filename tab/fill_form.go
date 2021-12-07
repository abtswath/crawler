package tab

import (
	"crawler/utils"
	"github.com/go-rod/rod"
	"strconv"
	"strings"
	"time"
)

type fillMethod func(element *rod.Element)

var (
	inputFillMethods = map[string]fillMethod{
		InputTypeSearch:        fillText,
		InputTypeNumber:        fillNumber,
		InputTypePassword:      fillInputByType,
		InputTypeEmail:         fillInputByType,
		InputTypeTel:           fillInputByType,
		InputTypeURL:           fillInputByType,
		InputTypeDate:          fillDate,
		InputTypeDatetimeLocal: fillDate,
		InputTypeMonth:         fillDate,
		InputTypeTime:          fillDate,
		InputTypeWeek:          fillDate,
		InputTypeCheckbox:      fillCheckbox,
		InputTypeRadio:         fillRadio,
	}
)

func (t *Tab) fillForm() {
	go t.input()
	go t.selectOption()
	go t.inputTextarea()
}

func (t *Tab) input() {
	defer t.wg.Done()
	elements, err := t.Page.ElementsByJS(rod.Eval(`document.querySelectorAll('input')`))
	if err != nil {
		t.logger.Debugf("Get tag input error: %s\n", err)
		return
	}
	for _, element := range elements {
		inputType := elementAttributeValue(element, "type", "text")
		if inputType == InputTypeFile {
			uploadFile(element, t.uploadFile)
			continue
		}
		if f, ok := inputFillMethods[inputType]; ok {
			f(element)
		}
	}
}

func (t *Tab) selectOption() {
	defer t.wg.Done()
	elements, err := t.Page.ElementsByJS(rod.Eval(`document.querySelectorAll('select')`))
	if err != nil {
		t.logger.Debugf("Get tag select error: %s\n", err)
		return
	}
	for _, element := range elements {
		_ = element.Select([]string{`nth-of-type(2)`}, true, rod.SelectorTypeCSSSector)
	}
}

func (t *Tab) inputTextarea() {
	defer t.wg.Done()
	elements, err := t.Page.ElementsByJS(rod.Eval(`document.querySelectorAll('textarea')`))
	if err != nil {
		t.logger.Debugf("Get tag textarea error: %s\n", err)
		return
	}
	for _, element := range elements {
		_ = element.Input(getValidInputTextValue(element))
	}
}

func fillText(element *rod.Element) {
	inputName := elementAttributeValue(element, "name", "")
	if inputName == "" {
		_ = element.Input(getValidInputTextValue(element))
		return
	}
	for _, item := range PredictableInputValues {
		if utils.StringArrayInclude(item.Keyword, inputName) {
			_ = element.Input(item.Value)
			return
		}
	}
	_ = element.Input(getValidInputTextValue(element))
}

func fillNumber(element *rod.Element) {
	minValue := elementAttributeValue(element, "min", "")
	maxValue := elementAttributeValue(element, "max", "")
	if minValue != "" {
		_ = element.Input(minValue)
		return
	}
	if maxValue != "" {
		_ = element.Input(maxValue)
		return
	}
	_ = element.Input(PredictableInputValues["number"].Value)
}

func fillInputByType(element *rod.Element) {
	inputType := elementAttributeValue(element, "type", "text")
	if inputType != "" {
		return
	}
	_ = element.Input(PredictableInputValues[inputType].Value)
}

func fillDate(element *rod.Element) {
	if InputTypeWeek == elementAttributeValue(element, "type", "text") {
		_ = element.Input("1")
		return
	}
	_ = element.InputTime(time.Now())
}

func fillCheckbox(element *rod.Element) {
	checked, err := element.Property("checked")
	if err != nil {
		return
	}
	if !checked.Bool() {
		element.MustClick()
	}
}

func fillRadio(element *rod.Element) {

}

func uploadFile(element *rod.Element, file ...string) {
	_ = element.SetFiles(file)
}

func elementAttributeValue(element *rod.Element, attribute string, defaultValue string) string {
	attributeValue, _ := element.Attribute(attribute)
	if attributeValue == nil {
		return defaultValue
	}
	return strings.ToLower(*attributeValue)
}

func getValidInputTextValue(element *rod.Element) string {
	length := 10
	var err error
	minLengthValue := elementAttributeValue(element, "minlength", "")
	maxLengthValue := elementAttributeValue(element, "maxlength", "")
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
