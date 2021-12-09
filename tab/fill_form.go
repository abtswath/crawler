package tab

import (
	"crawler/utils"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"strconv"
	"strings"
	"sync"
	"time"
)

type inputMethod func(element *rod.Element)

type form struct {
	form         *rod.Element
	tab          *Tab
	file         string
	wg           sync.WaitGroup
	inputMethods map[string]inputMethod
}

func (t *Tab) formSubmit() {
	defer t.wg.Done()
	forms, err := t.Page.Elements("form")
	if err != nil {
		return
	}
	for _, element := range forms {
		if !isVisibility(element) {
			continue
		}
		_, err = element.Eval(`this.target = '_self'`)
		if err != nil {
			t.logger.Debugf("Modify form [%s] target error: %s", element.String(), err)
		}
		f := form{
			form: element,
			tab:  t,
			file: t.uploadFile,
		}
		f.inputMethods = map[string]inputMethod{
			InputTypeSearch:        f.inputText,
			InputTypeText:          f.inputText,
			InputTypeNumber:        f.inputNumber,
			InputTypePassword:      f.inputByType,
			InputTypeEmail:         f.inputByType,
			InputTypeTel:           f.inputByType,
			InputTypeURL:           f.inputByType,
			InputTypeDate:          f.inputDate,
			InputTypeDatetimeLocal: f.inputDate,
			InputTypeMonth:         f.inputDate,
			InputTypeTime:          f.inputDate,
			InputTypeWeek:          f.inputDate,
			InputTypeCheckbox:      f.inputCheckbox,
			InputTypeRadio:         f.inputRadio,
		}
		f.wg.Add(3)
		f.fill()
		f.wg.Add(2)
		f.submit()
		f.wg.Wait()
	}
}

func (f *form) fill() {
	go f.input()
	go f.selectOption()
	go f.inputTextarea()
}

func (f *form) submit() {
	go f.attachSubmitEvent()
	go f.clickButtons()
}

func (f *form) attachSubmitEvent() {
	defer f.wg.Done()
	_, err := f.form.Eval(`this.submit()`)
	if err != nil {
		f.tab.logger.Debugf("Form [%s] submit error: %s", f.form.String(), err)
	}
}

func (f *form) clickButtons() {
	defer f.wg.Done()
	submitButtons, err := f.form.ElementsByJS(rod.Eval(`document.querySelectorAll('input[type=submit]')`))
	if err != nil {
		f.tab.logger.Debugf("Get submit buttons of form error: %s", err)
	} else {
		for _, button := range submitButtons {
			if elementAttributeValue(button, "type", "") == "reset" {
				continue
			}
			if !isVisibility(button) {
				continue
			}
			err = button.Click(proto.InputMouseButtonLeft)
			if err != nil {
				f.tab.logger.Debugf("Input submit [%s] clicked error: %s", button.String(), err)
			}
		}
	}
	buttons, err := f.form.Elements("button")
	if err != nil {
		f.tab.logger.Debugf("Get buttons of form error: %s", err)
		return
	}
	for _, button := range buttons {
		if elementAttributeValue(button, "type", "") == "reset" {
			continue
		}
		if !isVisibility(button) {
			continue
		}
		err = button.Click(proto.InputMouseButtonLeft)
		if err != nil {
			f.tab.logger.Debugf("Button [%s] clicked error: %s", button.String(), err)
		}
	}
}

func (f *form) input() {
	defer f.wg.Done()
	elements, err := f.form.ElementsByJS(rod.Eval(`document.querySelectorAll('input')`))
	if err != nil {
		f.tab.logger.Debugf("Get tag input error: %s\n", err)
		return
	}
	for _, element := range elements {
		inputType := elementAttributeValue(element, "type", "text")
		if inputType == InputTypeFile {
			f.setFiles(element, f.file)
			continue
		}
		if m, ok := f.inputMethods[inputType]; ok {
			m(element)
		} else {
			f.inputText(element)
		}
	}
}

func (f *form) selectOption() {
	defer f.wg.Done()
	elements, err := f.form.ElementsByJS(rod.Eval(`document.querySelectorAll('select')`))
	if err != nil {
		f.tab.logger.Debugf("Get tag select error: %s\n", err)
		return
	}
	for _, element := range elements {
		if !isVisibility(element) {
			continue
		}
		err = element.Select([]string{`:first-child`}, true, rod.SelectorTypeCSSSector)
		if err != nil {
			f.tab.logger.Debugf("Select selection error: %s", err)
		}
	}
}

func (f *form) inputTextarea() {
	defer f.wg.Done()
	elements, err := f.form.ElementsByJS(rod.Eval(`document.querySelectorAll('textarea')`))
	if err != nil {
		f.tab.logger.Debugf("Get tag textarea error: %s\n", err)
		return
	}
	for _, element := range elements {
		err = element.Input(getValidInputTextValue(element))
		if err != nil {
			f.tab.logger.Debugf("Input textarea error: %s", err)
		}
	}
}

func (f *form) inputText(element *rod.Element) {
	inputName := elementAttributeValue(element, "name", "")
	if inputName == "" {
		err := element.Input(getValidInputTextValue(element))
		if err != nil {
			f.tab.logger.Debugf("Input text error: %s", err)
		}
		return
	}
	for _, item := range PredictableInputValues {
		if utils.StringArrayInclude(item.Keyword, inputName) {
			err := element.Input(item.Value)
			if err != nil {
				f.tab.logger.Debugf("Input text [%s] error: %s", inputName, err)
			}
			return
		}
	}
	err := element.Input(getValidInputTextValue(element))
	if err != nil {
		f.tab.logger.Debugf("Input text [%s] error: %s", inputName, err)
	}
}

func (f *form) inputNumber(element *rod.Element) {
	minValue := elementAttributeValue(element, "min", "")
	maxValue := elementAttributeValue(element, "max", "")
	if minValue != "" {
		err := element.Input(minValue)
		if err != nil {
			f.tab.logger.Debugf("Input number error: %s", err)
		}
		return
	}
	if maxValue != "" {
		err := element.Input(maxValue)
		if err != nil {
			f.tab.logger.Debugf("Input number error: %s", err)
		}
		return
	}
	err := element.Input(PredictableInputValues["number"].Value)
	if err != nil {
		f.tab.logger.Debugf("Input number error: %s", err)
	}
}

func (f *form) inputByType(element *rod.Element) {
	inputType := elementAttributeValue(element, "type", "text")
	if inputType == "" {
		return
	}
	err := element.Input(PredictableInputValues[inputType].Value)
	if err != nil {
		f.tab.logger.Debugf("Input [%s] error: %s", inputType, err)
	}
}

func (f *form) inputDate(element *rod.Element) {
	if InputTypeWeek == elementAttributeValue(element, "type", "text") {
		err := element.Input("1")
		if err != nil {
			f.tab.logger.Debugf("Input week error: %s", err)
		}
		return
	}
	err := element.InputTime(time.Now())
	if err != nil {
		f.tab.logger.Debugf("Input time error: %s", err)
	}
}

func (f *form) inputCheckbox(element *rod.Element) {
	checked, err := element.Property("checked")
	if err != nil {
		return
	}
	if !checked.Bool() {
		err = element.Click(proto.InputMouseButtonLeft)
		if err != nil {
			fmt.Printf("Click error: %s\n", err)
		}
	}
}

func (f *form) inputRadio(element *rod.Element) {

}

func (f *form) setFiles(element *rod.Element, file ...string) {
	_, err := element.Eval(`this.removeAttribute('accept')`)
	if err != nil {
		f.tab.logger.Debugf("Input file [%s] remove attribute accept error: %s", element.String(), err)
	}
	err = element.SetFiles(file)
	if err != nil {
		f.tab.logger.Debugf("Set files error: %s", err)
	}
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

func isVisibility(element *rod.Element) bool {
	visibility, err := element.Visible()
	if err != nil {
		return false
	}
	return visibility
}
