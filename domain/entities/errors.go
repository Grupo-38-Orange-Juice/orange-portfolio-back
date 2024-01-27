package entities

import "fmt"

func EmptyFieldError(field string) error {
	return fmt.Errorf("\"%s\" vazio", field)
}

func InvalidFieldPatternError(field string) error {
	return fmt.Errorf("\"%s\" com padrão inválido", field)
}

func InvalidFieldLengthError(field string, min int, max int) error {
	return fmt.Errorf("\"%s\" tamanho inválido, o campo deve ter entre min: %d, max: %d caracteres", field, min, max)
}

func NegativeIntegerError(field string) error {
	return fmt.Errorf("\"%s\" esse valor não pode ser negativo", field)
}

func InvalidEmailError() error {
	return fmt.Errorf("formato de email inválido")
}

func InvalidPatternError(field, rule string) error {
	return fmt.Errorf("\"%s\" deve ter pelo menos um  %s carácter", field, rule)
}

func InvalidFieldRangeError(field string, min int, max int) error {
	return fmt.Errorf("\"%s\" fora do intervalo, o campo deve ter entre min: %d, max: %d", field, min, max)
}
