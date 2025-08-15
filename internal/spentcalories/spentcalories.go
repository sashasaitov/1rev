package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректный формат данных: ожидается три элемента, получено %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при преобразовании количества шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге продолжительности: %v", err)
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	steplength := height * stepLengthCoefficient
	totaldistance := steplength * float64(steps)
	inkm := totaldistance / mInKm
	return inkm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	hours := duration.Hours()
	speed := dist / hours
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return "", err
	}

	speed := meanSpeed(steps, height, duration)
	if speed < 0 {
		return "", fmt.Errorf("ошибка при расчёте скорости")
	}

	distance := distance(steps, height)

	caloriesr, err := RunningSpentCalories(steps, weight, height, duration)
	caloriesw, err := WalkingSpentCalories(steps, weight, height, duration)

	var trainingStr string
	switch strings.ToLower(activityType) {
	case "бег":
		trainingStr = fmt.Sprintf(
			`Тип тренировки: Бег
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f`,
			duration.Hours(),
			distance,
			speed,
			caloriesr,
		)
	case "ходьба":
		trainingStr = fmt.Sprintf(
			`Тип тренировки: Ходьба
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f`,
			duration.Hours(),
			distance,
			speed,
			caloriesw,
		)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return trainingStr, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("количество шагов не может быть отрицательным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	speed := meanSpeed(steps, height, duration)
	if speed < 0 {
		return 0, fmt.Errorf("Ошибка при расчете скорости")
	}
	minutes := duration.Minutes()
	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, fmt.Errorf("количество шагов не может быть отрицательным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	speed := meanSpeed(steps, height, duration)
	if speed < 0 {
		return 0, fmt.Errorf("Ошибка при расчете скорости")
	}
	minutes := duration.Minutes()
	calories := ((weight * minutes * speed) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
