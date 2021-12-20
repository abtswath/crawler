package tab

import (
	"crawler/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"strconv"
	"strings"
	"sync"
	"time"
)

type inputMethod func(element *rod.Element) error

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
		f.fill()
		f.submit()
	}
}

func (f *form) fill() {
	f.wg.Add(1)
	go f.input()
	f.wg.Add(1)
	go f.selectOption()
	f.wg.Add(1)
	go f.inputTextarea()
	f.wg.Wait()
}

func (f *form) submit() {
	f.wg.Add(1)
	go f.attachSubmitEvent()
	f.wg.Add(1)
	go f.clickButtons()
	f.wg.Wait()
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
			if !isVisibility(button) {
				continue
			}
			if elementAttributeValue(button, "type", "") == "reset" {
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
		if !isVisibility(button) {
			continue
		}
		if elementAttributeValue(button, "type", "") == "reset" {
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
		if !isVisibility(element) {
			continue
		}
		inputType := elementAttributeValue(element, "type", "text")
		if inputType == InputTypeHidden {
			continue
		}
		if inputType == InputTypeFile {
			err = f.setFiles(element, f.file)
			if err != nil {
				f.tab.logger.Debugf("Element %s set file error: %s", element.String(), err)
			}
			continue
		}
		if m, ok := f.inputMethods[inputType]; ok {
			err = m(element)
			if err != nil {
				f.tab.logger.Debugf("Element %s input error: %s", element.String(), err)
			}
		} else {
			err = f.inputText(element)
			if err != nil {
				f.tab.logger.Debugf("Element %s input error: %s", element.String(), err)
			}
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
		if !isVisibility(element) {
			continue
		}
		err = element.Input(getValidInputTextValue(element))
		if err != nil {
			f.tab.logger.Debugf("Input textarea error: %s", err)
		}
	}
}

func (f *form) inputText(element *rod.Element) error {
	inputName := elementAttributeValue(element, "name", "")
	if inputName == "" {
		err := element.Input(getValidInputTextValue(element))
		if err != nil {
			f.tab.logger.Debugf("Input text error: %s", err)
			return err
		}
		return nil
	}
	for _, item := range PredictableInputValues {
		if utils.StringArrayInclude(item.Keyword, inputName) {
			return element.Input(item.Value)
		}
	}
	return element.Input(getValidInputTextValue(element))
}

func (f *form) inputNumber(element *rod.Element) error {
	minValue := elementAttributeValue(element, "min", "")
	maxValue := elementAttributeValue(element, "max", "")
	if minValue != "" {
		return element.Input(minValue)
	}
	if maxValue != "" {
		return element.Input(maxValue)
	}
	return element.Input(PredictableInputValues["number"].Value)
}

func (f *form) inputByType(element *rod.Element) error {
	inputType := elementAttributeValue(element, "type", "text")
	if inputType == "" {
		return element.Input(getValidInputTextValue(element))
	}
	err := element.Input(PredictableInputValues[inputType].Value)
	if err != nil {
	}
	return nil
}

func (f *form) inputDate(element *rod.Element) error {
	if InputTypeWeek == elementAttributeValue(element, "type", "text") {
		return element.Input("1")
	}
	return element.InputTime(time.Now())
}

func (f *form) inputCheckbox(element *rod.Element) error {
	checked, err := element.Property("checked")
	if err != nil {
		return err
	}
	if !checked.Bool() {
		return element.Click(proto.InputMouseButtonLeft)
	}
	return nil
}

func (f *form) inputRadio(element *rod.Element) error {
	return nil
}

func (f *form) setFiles(element *rod.Element, file ...string) error {
	_, err := element.Eval(`this.removeAttribute('accept')`)
	if err != nil {
		return err
	}
	return element.SetFiles(file)
}

func elementAttributeValue(element *rod.Element, attribute string, defaultValue string) string {
	attributeValue, _ := element.Attribute(attribute)
	if attributeValue == nil {
		return defaultValue
	}
	value := *attributeValue
	return strings.ToLower(value)
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
