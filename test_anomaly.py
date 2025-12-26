import requests
import random
import time

url = "http://localhost:8080/data"

print("1. Отправляем 60 нормальных значений (обучаем окно)...")
for i in range(60):
    # Генерируем "нормальную" нагрузку: число от 10 до 20
    val = random.uniform(10.0, 20.0)
    payload = {"timestamp": int(time.time()), "value": val}
    try:
        requests.post(url, json=payload)
    except:
        pass
    # Небольшая пауза, чтобы не забить канал мгновенно, хотя Go справится
    time.sleep(0.05)

print("2. Окно заполнено. Отправляем АНОМАЛИЮ (значение 100)...")
requests.post(url, json={"timestamp": int(time.time()), "value": 100.0})

print("Готово! Смотри логи Go-сервиса.")
