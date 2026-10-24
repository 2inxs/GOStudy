/*
Тут я тоже честно говоря, не понял, что требуется.
Структуры я уже поверхностно пробовал, поэтому вот реализация

Ответ на вопрос:
Для частичного обновления нам нужно 3 состояния, а bool дает только true/false
При обычном Newsletter bool, при создании функции компилятор названичит false
Если же мы создаем с указателем, то явно передаем nil
*/
package main

type Profile struct {
	Name       string
	Age        int
	Email      string
	Newsletter bool
}

type ProfilePatch struct {
	Name       *string
	Age        *int
	Email      *string
	Newsletter *bool
}

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if p == nil {
		return
	}

	if patch.Name != nil {
		p.Name = *patch.Name
	}

	if patch.Age != nil {
		p.Age = *patch.Age
	}

	if patch.Email != nil {
		p.Email = *patch.Email
	}

	if patch.Newsletter != nil {
		p.Newsletter = *patch.Newsletter
	}
}
