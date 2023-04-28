package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	outputFilePath = "./annotationsOutput"
	configFile     = "./remapConfig"
)

func reMapper(filePath string, fileName string, remapConfig map[int]int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// создаем новый файл для записи измененных строк
	newFile, err := os.Create(outputFilePath + "/" + fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()
	// создаем сканер для чтения строк из файла
	scanner := bufio.NewScanner(file)
	// проходим по каждой строке файла
	for scanner.Scan() {
		// разбиваем строку на части по пробелам
		parts := strings.Split(scanner.Text(), " ")

		// проверяем, что первая часть строки может быть преобразована в число
		num, err := strconv.Atoi(parts[0])
		if err != nil {
			// если не может, то просто записываем строку в новый файл без изменений
			return err
		}
		//Если значение есть в конфигурации, проведем изменение
		if val, ok := remapConfig[num]; ok {
			num = val
		}
		// объединяем измененное число со всеми остальными частями строки
		newLine := strconv.Itoa(num) + " " + strings.Join(parts[1:], " ")
		// записываем измененную строку в новый файл
		fmt.Fprintln(newFile, newLine)

	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func readConfig() (map[int]int, error) {
	resConfig := make(map[int]int, 30)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// создаем сканер для чтения строк из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
		// разбиваем строку на части по пробелам
		parts := strings.Split(scanner.Text(), " ")
		// проверяем, что первая часть строки может быть преобразована в число
		orig, err := strconv.Atoi(parts[0])
		if err != nil {
			// если не может, то просто записываем строку в новый файл без изменений
			log.Println("тут")
			return nil, err
		}
		desired, err := strconv.Atoi(parts[1])
		if err != nil {
			// если не может, то просто записываем строку в новый файл без изменений
			log.Println("здесь")
			return nil, err
		}
		resConfig[orig] = desired
	}
	return resConfig, nil
}

func main() {
	//Чтение файла конфигурации
	remapConfig := make(map[int]int, 30)
	remapConfig, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}
	// specify the original directory path
	dirPath := "./annotationsInput"
	// specify the new directory path
	newDirPath := "./annotationsOutput"
	// create the new directory if it doesn't exist
	err = os.MkdirAll(newDirPath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	// read all files in the original directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	// loop through each file and modify its contents
	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			if strings.HasSuffix(file.Name(), ".txt") {
				filePath := dirPath + "/" + file.Name()
				err = reMapper(filePath, file.Name(), remapConfig)
				if err != nil {
					fmt.Println(err)
				}
				defer wg.Done()
			}
		}(file)
	}
	wg.Wait()

}
