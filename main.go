package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"uvm-assembler/assembler"
)

func main() {

	inputFile := flag.String("input", "", "Путь к исходному файлу с текстом программы")
	outputFile := flag.String("output", "", "Путь к двоичному файлу-результату")
	testMode := flag.Bool("test", false, "Режим тестирования (вывод промежуточного представления)")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Необходимо указать входной файл")
		fmt.Println("Использование: uvm-assembler -input program.asm [-output program.bin] [-test]")

		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Файл %s не найден\n", *inputFile)
		os.Exit(1)
	}

	if *outputFile == "" {
		fmt.Println("Необходимо указать файл-результата")
		fmt.Println("Использование: uvm-assembler [-input program.asm] -output program.bin [-test]")

		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Файл %s не найден\n", *inputFile)
		os.Exit(1)
	}

	fmt.Println("===== Ассемблер УВМ =====")
	fmt.Println("=======================================")
	fmt.Printf("Входной файл:  %s\n", *inputFile)
	fmt.Printf("Выходной файл: %s\n", *outputFile)
	fmt.Printf("Режим тестирования: %v\n", *testMode)
	fmt.Println()

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("❌ Ошибка чтения файла: %v\n", *inputFile)
		os.Exit(1)
	}

	fmt.Printf("✅ Файл прочитан успешно (%d байт)\n", len(content))

	parser := assembler.NewParser(string(content))
	commands, err := parser.Parse()
	if err != nil {
		fmt.Printf("Ошибка парсинга: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Программа разобрана успешно (%d команд)\n", len(commands))

	if *testMode {
		displayTestResults(commands)
	}
	encoder := assembler.NewEncoder()
	binaryProgram := make([]byte, 0, len(commands)*5) 
	
	if *outputFile != "" {
		for i, cmd := range commands {
		machineCode, err := encoder.Encode(cmd)
		if err != nil {
			fmt.Printf("❌ Ошибка кодирования команды %d: %v\n", i+1, err)
			os.Exit(1)
    	}

		binaryProgram = append(binaryProgram, machineCode...)
		fmt.Printf("✅ Команда %d закодирована: %s\n", 
			i+1, encoder.BytesToHexString(machineCode))
		}
	}

	err = os.WriteFile(*outputFile, binaryProgram, 0644)
	if err != nil {
		fmt.Printf("❌ Ошибка записи файла: %v\n", err)
        os.Exit(1)
	}

	fileInfo, _ := os.Stat(*outputFile)
	fmt.Printf("\n💾 Размер двоичного файла: %d байт\n", fileInfo.Size())
    fmt.Printf("📦 Количество команд: %d\n", len(commands))
    fmt.Printf("💿 Общий размер: %d байт (%d команд × 5 байт)\n", 
        len(commands)*5, len(commands))

	if *testMode {
        fmt.Println("\n БАЙТОВОЕ ПРЕДСТАВЛЕНИЕ (как в спецификации):")
        fmt.Println("==============================================")
        
        for i, cmd := range commands {
            machineCode, _ := encoder.Encode(cmd)
            fmt.Printf("Команда %d: %s\n", i+1, encoder.BytesToHexString(machineCode))
        }
        
        fmt.Println("\n СРАВНЕНИЕ С ТЕСТАМИ ИЗ СПЕЦИФИКАЦИИ:")
        fmt.Println("======================================")
        verifyByteTests(commands, encoder)
    }
}

// displayTestResults выводит результаты в формате как в спецификации УВМ
func displayTestResults(commands []assembler.Command) {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("🔍 РЕЖИМ ТЕСТИРОВАНИЯ - ПРОМЕЖУТОЧНОЕ ПРЕДСТАВЛЕНИЕ")
	fmt.Println(strings.Repeat("═", 60))
	
	for i, cmd := range commands {
		fmt.Printf("\nКоманда %d:\n", i+1)
		fmt.Printf("  Мнемоника: %s\n", cmd.Type.TypeName())
		fmt.Printf("  Поля: %s\n", cmd.ToTestFormat())
		fmt.Printf("  Детали:\n")
		
		for field, value := range cmd.Fields {
			fmt.Printf("    %s: %d\n", field, value)
		}
	}
	
	// 🧪 ПУНКТ 6: Проверка тестовых случаев из спецификации
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("🧪 ПРОВЕРКА ТЕСТОВЫХ СЛУЧАЕВ ИЗ СПЕЦИФИКАЦИИ УВМ")
	fmt.Println(strings.Repeat("═", 60))
	
	verifySpecificationTests(commands)
}

// verifySpecificationTests проверяет соответствие тестам из спецификации
func verifySpecificationTests(commands []assembler.Command) {
	// Ожидаемые результаты из спецификации УВМ
	expectedTests := []struct {
		name     string
		expected map[string]uint32
	}{
		{
			"Загрузка константы",
			map[string]uint32{"A": 59, "B": 9, "C": 771},
		},
		{
			"Чтение значения из памяти", 
			map[string]uint32{"A": 8, "B": 499, "C": 42, "D": 35},
		},
		{
			"Запись значения в память",
			map[string]uint32{"A": 37, "B": 25, "C": 3},
		},
		{
			"Унарная операция: sqrt()",
			map[string]uint32{"A": 4, "B": 9, "C": 804},
		},
	}
	
	allTestsPassed := true
	
	for i, test := range expectedTests {
		fmt.Printf("\nТест %d: %s\n", i+1, test.name)
		fmt.Printf("  Ожидается: %v\n", formatExpected(test.expected))
		
		if i < len(commands) {
			cmd := commands[i]
			fmt.Printf("  Получено:  %s\n", cmd.ToTestFormat())
			
			// Проверяем соответствие полей
			testPassed := true
			for field, expectedValue := range test.expected {
				actualValue, exists := cmd.Fields[field]
				if !exists || actualValue != expectedValue {
					testPassed = false
					allTestsPassed = false
					fmt.Printf("  ❌ Поле %s: ожидалось=%d, получено=%d\n", 
						field, expectedValue, actualValue)
				}
			}
			
			if testPassed {
				fmt.Printf("  ✅ Тест пройден!\n")
			} else {
				fmt.Printf("  ❌ Тест не пройден!\n")
			}
		} else {
			fmt.Printf("  ❌ Нет команды для теста!\n")
			allTestsPassed = false
		}
	}
	
	fmt.Println("\n" + strings.Repeat("═", 60))
	if allTestsPassed {
		fmt.Println("🎉 ВСЕ ТЕСТЫ ИЗ СПЕЦИФИКАЦИИ ПРОЙДЕНЫ УСПЕШНО!")
	} else {
		fmt.Println("💥 НЕКОТОРЫЕ ТЕСТЫ НЕ ПРОЙДЕНЫ!")
	}
	fmt.Println(strings.Repeat("═", 60))
}

// formatExpected форматирует ожидаемые значения для красивого вывода
func formatExpected(expected map[string]uint32) string {
	if len(expected) == 3 {
		return fmt.Sprintf("(A=%d, B=%d, C=%d)", expected["A"], expected["B"], expected["C"])
	} else if len(expected) == 4 {
		return fmt.Sprintf("(A=%d, B=%d, C=%d, D=%d)", expected["A"], expected["B"], expected["C"], expected["D"])
	}
	return fmt.Sprintf("%v", expected)
}

// verifyByteTests проверяет соответствие байтовым тестам из спецификации
func verifyByteTests(commands []assembler.Command, encoder *assembler.Encoder) {
	expectedByteTests := []struct {
		name     string
		expected []byte
	}{
		{
			"Загрузка константы (A=59, B=9, C=771)",
			[]byte{0x7B, 0x32, 0x30, 0x00, 0x00},
		},
		{
			"Чтение из памяти (A=8, B=499, C=42, D=35)",
			[]byte{0xC8, 0x7C, 0x80, 0x3A, 0x02},
		},
		{
			"Запись в память (A=37, B=25, C=3)", 
			[]byte{0x65, 0x36, 0x00, 0x00, 0x00},
		},
		{
			"Квадратный корень (A=4, B=9, C=804)",
			[]byte{0x44, 0x42, 0x32, 0x00, 0x00},
		},
	}

	allTestsPassed := true

	for i, test := range expectedByteTests {
		fmt.Printf("\nТест %d: %s\n", i+1, test.name)
		fmt.Printf("  Ожидается: %s\n", encoder.BytesToHexString(test.expected))
		
		if i < len(commands) {
			actual, err := encoder.Encode(commands[i])
			if err != nil {
				fmt.Printf("  ❌ Ошибка кодирования: %v\n", err)
				allTestsPassed = false
				continue
			}
			
			fmt.Printf("  Получено:  %s\n", encoder.BytesToHexString(actual))
			
			// Сравниваем байты
			match := true
			for j := range test.expected {
				if test.expected[j] != actual[j] {
					match = false
					break
				}
			}
			
			if match {
				fmt.Printf("  ✅ Байты совпадают!\n")
			} else {
				fmt.Printf("  ❌ Байты не совпадают!\n")
				allTestsPassed = false
			}
		} else {
			fmt.Printf("  ❌ Нет команды для теста!\n")
			allTestsPassed = false
		}
	}
	
	fmt.Println("\n" + strings.Repeat("═", 60))
	if allTestsPassed {
		fmt.Println("🎉 ВСЕ БАЙТОВЫЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!")
	} else {
		fmt.Println("💥 НЕКОТОРЫЕ БАЙТОВЫЕ ТЕСТЫ НЕ ПРОЙДЕНЫ!")
	}
	fmt.Println(strings.Repeat("═", 60))
}