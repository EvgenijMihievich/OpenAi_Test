package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

func main() {
	// Укажите ваш API-ключ OpenAI
	apiKey := "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"

	// Создаем буфер для чтения ввода пользователя из консоли
	reader := bufio.NewReader(os.Stdin)

	// Создаем пустой массив для сохранения истории диалога
	var userMessages []ChatMessage

	// Бесконечный цикл для чтения ввода пользователя и отправки запросов
	for {
		fmt.Print("Введите запрос (для выхода введите 'exit'): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении ввода:", err)
			return
		}

		// Убираем лишние пробелы и символы переноса строки
		input = strings.TrimSpace(input)

		// Проверяем, если пользователь хочет выйти
		if input == "exit" {
			fmt.Println("Программа завершена.")
			return
		}

		// Добавляем ввод пользователя в историю диалога
		userMessage := ChatMessage{Role: "user", Content: input}
		userMessages = append(userMessages, userMessage)

		// Проверяем, если длина истории диалога достигла максимального значения
		if len(userMessages) > 100 {
			userMessages = userMessages[len(userMessages)-100:]
		}

		// Создаем JSON-запрос
		requestBody, err := json.Marshal(ChatRequest{
			Model:     "gpt-3.5-turbo-1106",
			Messages:  userMessages,
			MaxTokens: 2048,
		})

		if err != nil {
			fmt.Println("Ошибка при создании JSON-запроса:", err)
			return
		}

		// Отправляем POST-запрос к OpenAI API с токеном авторизации
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))

		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

		// Считываем ответ
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка при чтении ответа:", err)
			return
		}

		// Выводим ответ
		fmt.Println("Ответ:")
		fmt.Println(string(responseBody))

		// Добавляем ответ системы в историю диалога
		userMessage = ChatMessage{}
		err = json.Unmarshal(responseBody, &userMessage)
		if err != nil {
			fmt.Println("Ошибка при разборе ответа:", err)
			return
		}
		userMessage.Role = "assistant"
		userMessages = append(userMessages, userMessage)
	}
}
