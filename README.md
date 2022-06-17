# Loki y Grafana

## Instalación
* Tanka (recomendado): Herramienta propia de grafana, simil a helm.
* Helm: Para entornos k8s.
* Docker (& compose)
* Binario localmente
* Código fuente

### Docker
```shell
docker run --name loki -p 3100:3100 grafana/loki:2.4.2
docker run --name promtail -link loki grafana/promtail:2.4.2
```
_promtail_: Herramienta para inyectar los logs.
*Importante*: Se necesita generar un datasource en provisioning/datasources para conectar a la fuente de los logs.
*Compose*: Se puede descargar del repositorio de grafana.
*flog*: Generan logs de forma aleatoria en formato json.

### Grafana Cloud (versión gratuita)
[Planes de Grafana](grafana.com/pricing)

## Clientes de logs para Grafana Loki
* Promtail
* Loki push API
* Docker driver
* Fluentd & Fuelt bit
* Logtash

## Estructura de log + labels
1. Etiquetas -streams-:
   1. Pares clave-valor; similar a Prometheus ej: `{cluster="cluster-01", instance="instance-02"}`
2. Líneas de log:
   1. Pares de fechas y mensaje
   2. Ordenadas cronológicamente*

## Escribiendo logs

### Enviando logs via Loki Push API
Enviando JSON HTTP API
Request: 
```shell
curl --location --request POST 'localhost:3100/loki/api/v1/push' \
--header 'Content-Type: application/json' \
--data-raw '{
    "streams": [
        {
            "stream": {
                "cluster": "cluster-01",
                "instance": "instance-01"
            },
            "values": [
                ["1653414480164000000","Hola mundo, soy un log"]
            ]
        }
    ]
}'
```

### Enviando desde Logstash
[Plugin para logstash](https://grafana.com/docs/loki/latest/clients/logstash/)

### Logs desde el Standard output de Docker
Se utliza un plugin:
```shell
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
```
## Busqueda en Loki
Utiliza LogQL
Para buscar se utiliza el {} y se busca dentro los campos.

### Pipelines
Una query de LogQL se divide en 2 elementos:
* ````Log Stream selector + Log Pipeline````
A su vez tienen 3 expresiones
1. Filtros por lineas
2. Parseo
3. Formateo

#### Filtros
Por linea: 
{job="nginx"} |= "error" // Que contenga
{job="nginx"} |~ "error=\w+" // Que contenga con regex
Por etiqueta que requieren "parseo" previo:
{job="nginx"} | duration > 1m and bytes_consumed > 20MB

#### Parser
1. JSON: {job="nginx"} | json ----- {job="nginx"} | json first_queries="queries[0]"
2. Logfmt

#### Formateo
* Formateo de etiquetas: Cambiar valor de etiquetas

#### Example
```shell
{job="flogs"} | json | method="GET" | line_format "{{.method}} => {{.request}}"
```

### Queries para métricas
1. Log range aggregations:
   1. `count_over_time(<criterio>[5m])`
   2. `sum by (host) (rate(<criterio> [1m]))`
2. Unwrapped range aggregations:
   1. `sum by (host) (sum_over_time(<criterio>[1m]))`

### Built-in functions
1. Transformaciones
   1. `duration_seconds(<duration>)`
   2. `bytes(<bytes>)`
2. Operaciones matemáticas (agregación):
   1. `topk(10,sum(rate(<criterio>)) by (host))`

## Alertas en Loki

### Loki ruler
Alertmanager.
Compatible con prometheus. Al ser independiente, se puede implementar en Grafana.
Reglas:
1. Alerting rules: Definir alerta a partir de expresiones
2. Recording rules: Generar métricas en base a los logs.