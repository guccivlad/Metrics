# Metrics

В `deploy.sh` я ждал загрузки `Prometheus` с помощью `sleep 900`. Если можно как то умнее, то будет круо, если скажешь мне об этом)

Метрики можно посмотреть здесь - `http://http://localhost:9090`

Список метрик:
`log_calls_total`, `succes_log_calls`, `failed_log_calls`, `log_request_duration_seconds`, `log_request_duration_seconds_avg`

Для последней метрики написан `PrometheusRule`
