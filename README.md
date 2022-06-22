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

## Prácticas recomendadas

### Truquitos para la ingestión de logs

Para ello tenemos que comprender el significado de:
* chunk_target_size, que indica el tamaño (en bytes) deseado por chunk una vez comprimido.
* max_chunk_age, que indica la máxima duración del chunk en memoria.
* chunk_idle_period, que indica la máxima duración del chunk en memoria, sin actualizaciones.

Una vez aprendido eso, algunas de las ideas serían:

* Usar etiquetas estáticas y/o de baja cardinalidad.
* Ajustar y refinar el chunk_target_size, la max_chunk_age y el chunk_idle_period acorde con nuestras necesidades.
* Usar el flag --analyze-labels para identificar etiquetas problemáticas.

### Respetar el orden cronológico de los logs

Desde una de las versiones recientes de Loki es posible inyectarle líneas de logs desordenadas, es decir, sin respetar 
el orden cronológico. Sin embargo, aunque dicha opción esté activada por defecto, es recomendable intentar prescindir 
de ella e incluso desactivarla.

Algunas de las estrategias para respetar dicho orden cuándo no se cumple, son:
* Identificar nuevos streams, es decir, nuevas etiquetas.
* Delegar la asignación de la fecha al cliente (Promtail).
* Gestionarlo y controlarlo desde nuestra aplicación.

### Aplicar el concepto de dead man’s switch

Finalmente, otra práctica recomendada sería configurar alertas siguiendo dicho concepto por ejemplo, mediante el uso 
de funciones como absent_over_time, que nos permitirán no solo actuar cuándo se cumplan ciertas condiciones en 
nuestros logs sino también cuándo no se cumplan (por ejemplo, cuándo no nos llega ningún log).
