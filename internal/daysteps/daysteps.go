package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		log.Println("Ошибка:")
		return 0, 0, fmt.Errorf("некорректный формат данных: ожидается два элемента, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println("Ошибка:", err)
		return 0, 0, fmt.Errorf("ошибка при преобразовании количества шагов: %v", err)
	}
	if steps <= 0 {
		log.Println("Ошибка:")
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		log.Println("Ошибка:", err)
		return 0, 0, fmt.Errorf("ошибка при парсинге продолжительности: %v", err)
	}
	if duration <= 0 {
		log.Println("Ошибка:")
		return 0, 0, fmt.Errorf("неверная продолжительность")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}
	if steps <= 0 {
		return "некоректное количество шагов"
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

	return result
}
