package assembler

import "fmt"

type CommandType int

const (
	LOAD_CONST CommandType = 59
	READ_MEM   CommandType = 8
	WRITE_MEM  CommandType = 37
	SQRT_OP    CommandType = 4
)

type Command struct {
	Type   CommandType
	Fields map[string]uint32
	Line   int
}

func (ct CommandType) TypeName() string {
	names := map[CommandType]string{
		LOAD_CONST: "LOAD",
		READ_MEM:   "READ",
		WRITE_MEM:  "WRITE",
		SQRT_OP:    "SQRT",
	}
	return names[ct]
}

type Parser struct {
	Lines       []string
	currentLine int
	filename    string
}

//Для 1 этапа

// String возвращает строковое представление команды в формате полей
func (c Command) String() string {
	switch c.Type {
	case LOAD_CONST:
		return fmt.Sprintf("A=%d, B=%d, C=%d", c.Fields["A"], c.Fields["B"], c.Fields["C"])
	case READ_MEM:
		return fmt.Sprintf("A=%d, B=%d, C=%d, D=%d", c.Fields["A"], c.Fields["B"], c.Fields["C"], c.Fields["D"])
	case WRITE_MEM, SQRT_OP:
		return fmt.Sprintf("A=%d, B=%d, C=%d", c.Fields["A"], c.Fields["B"], c.Fields["C"])
	default:
		return "Неизвестная команда"
	}
}

// ToTestFormat возвращает представление в формате как в тестах спецификации
func (c Command) ToTestFormat() string {
	switch c.Type {
	case LOAD_CONST:
		return fmt.Sprintf("(A=%d, B=%d, C=%d)", c.Fields["A"], c.Fields["B"], c.Fields["C"])
	case READ_MEM:
		return fmt.Sprintf("(A=%d, B=%d, C=%d, D=%d)", c.Fields["A"], c.Fields["B"], c.Fields["C"], c.Fields["D"])
	case WRITE_MEM, SQRT_OP:
		return fmt.Sprintf("(A=%d, B=%d, C=%d)", c.Fields["A"], c.Fields["B"], c.Fields["C"])
	default:
		return "(Неизвестная команда)"
	}
}