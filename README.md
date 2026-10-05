# Clima por CEP

Serviço em Go que recebe um CEP, identifica a cidade pelo [ViaCEP](https://viacep.com.br) e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin usando a [WeatherAPI](https://www.weatherapi.com).

## URL no Cloud Run

https://clima-cep-581002297807.europe-west1.run.app

Exemplo:

```
curl  https://clima-cep-581002297807.europe-west1.run.app/weather/01001000
```

## Endpoint

`GET /weather/{cep}`

| Cenário | Status | Resposta |
|---|--------|---|
| Sucesso | 200    | `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.65}` |
| CEP com formato inválido | 400    | `invalid zipcode` |
| CEP não encontrado | 404    | `can not find zipcode` |

O CEP deve ter exatamente 8 dígitos numéricos, sem hífen.

## Variáveis de ambiente

| Variável | Descrição |
|---|---|
| `WEATHER_API_KEY` | Chave da WeatherAPI (obrigatória) |
| `PORT` | Porta HTTP (padrão `8080`) |

## Rodando localmente com Docker

Crie o arquivo `.env` a partir do exemplo e coloque sua chave da WeatherAPI:

```
cp .env.example .env
```

Suba a aplicação:

```
docker compose up --build app
```

Teste:

```
curl http://localhost:8080/weather/01001000
```

## Rodando os testes

Com Docker:

```
docker compose run --rm test
```

Ou com Go instalado:

```
go test ./...
```

## Deploy no Cloud Run

```
gcloud run deploy clima-cep \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua_chave
```
