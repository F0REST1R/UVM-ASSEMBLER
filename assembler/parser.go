package assembler

import (
	"fmt"
	"strconv"
	"strings"
)

// Создает новый парсер
func NewParser(source string) *Parser {
	lines := strings.Split(source, "\n")
	return &Parser{
		Lines:       lines,
		currentLine: 0,
		filename:    "source.asm",
	}
}

func (p *Parser) Parse() ([]Command, error) {
	var commands []Command

	for lineNum, line := range p.Lines {
		line = strings.TrimSpace(line)

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Удаляем комментарий в конце строки
		if idx := strings.Index(line, ";"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		// Разбираем команду
		cmd, err := p.parseLine(line, lineNum+1)
		if err != nil {
			return nil, fmt.Errorf("строка %d: %v", lineNum+1, err)
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func (p *Parser) parseLine(line string, lineNum int) (Command, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Command{}, fmt.Errorf("пустая команда")
	}

	mnemonic := strings.ToUpper(parts[0])
	args := parts[1:]

	switch mnemonic {
	case "LOAD":
		return p.parseLoad(args, lineNum)
	case "READ":
		return p.parseRead(args, lineNum)
	case "WRITE":
		return p.parseWrite(args, lineNum)
	case "SQRT":
		return p.parseSqrt(args, lineNum)
	default:
		return Command{}, fmt.Errorf("неизвестная команда: %s", mnemonic)
	}
}

func (p *Parser) parseLoad(args []string, lineNum int) (Command, error) {
	if len(args) != 2 {
		return Command{}, fmt.Errorf("LOAD требует два аргумента: регистр, константа")
	}

	regB, err := p.parseRegister(args[0])
	if err != nil {
		return Command{}, err
	}

	constC, err := p.parseNumber(args[1])
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: LOAD_CONST,
		Fields: map[string]uint32{
			"A": 59, // Код операции
			"B": regB,
			"C": constC,
		},
		Line: lineNum,
	}, nil
}

// parseRead разбирает команду READ
func (p *Parser) parseRead(args []string, lineNum int) (Command, error) {
	if len(args) != 3 {
		return Command{}, fmt.Errorf("READ требует 3 аргумента: регистр_результата, смещение, базовый_регистр")
	}

	regD, err := p.parseRegister(args[0])
	if err != nil {
		return Command{}, err
	}

	offsetB, err := p.parseNumber(args[1])
	if err != nil {
		return Command{}, err
	}

	regC, err := p.parseRegister(args[2])
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: READ_MEM,
		Fields: map[string]uint32{
			"A": 8, // Код операции
			"B": offsetB,
			"C": regC,
			"D": regD,
		},
		Line: lineNum,
	}, nil
}

// parseWrite разбирает команду WRITE
func (p *Parser) parseWrite(args []string, lineNum int) (Command, error) {
	if len(args) != 2 {
		return Command{}, fmt.Errorf("WRITE требует 2 аргумента: регистр_значения, регистр_адреса")
	}

	regB, err := p.parseRegister(args[0])
	if err != nil {
		return Command{}, err
	}

	regC, err := p.parseRegister(args[1])
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: WRITE_MEM,
		Fields: map[string]uint32{
			"A": 37, // Код операции
			"B": regB,
			"C": regC,
		},
		Line: lineNum,
	}, nil
}

// parseSqrt разбирает команду SQRT
func (p *Parser) parseSqrt(args []string, lineNum int) (Command, error) {
	if len(args) != 2 {
		return Command{}, fmt.Errorf("SQRT требует 2 аргумента: регистр_источника, адрес_результата")
	}

	regB, err := p.parseRegister(args[0])
	if err != nil {
		return Command{}, err
	}

	addrC, err := p.parseNumber(args[1])
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: SQRT_OP,
		Fields: map[string]uint32{
			"A": 4, // Код операции
			"B": regB,
			"C": addrC,
		},
		Line: lineNum,
	}, nil
}

// parseRegister разбирает регистр (формат R0 - R63)
func (p *Parser) parseRegister(s string) (uint32, error) {
	if len(s) < 2 || s[0] != 'R' {
		return 0, fmt.Errorf("неверный формат регистра: %s, ожидается R0-R63", s)
	}

	regNum, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, fmt.Errorf("неверный номер регистра: %s", s)
	}

	if regNum < 0 || regNum > 63 {
		return 0, fmt.Errorf("номер регистра должен быть от 0 до 63: %s", s)
	}

	return uint32(regNum), nil
}

// parseNumber разбирает число (десятичное или шестнадцатеричное)
func (p *Parser) parseNumber(s string) (uint32, error) {
	// Пробуем разобрать как десятичное число
	if val, err := strconv.Atoi(s); err == nil {
		return uint32(val), nil
	}

	// Пробуем разобрать как шестнадцатеричное число
	if strings.HasPrefix(s, "0x") {
		if val, err := strconv.ParseInt(s[2:], 16, 32); err == nil {
			return uint32(val), nil
		}
	}

	return 0, fmt.Errorf("неверный числовой формат: %s", s)
}
