## Executando a Aplicação Localmente com Docker Compose

```bash
docker-compose up --build
```
4. A aplicação estará acessível localemnte em http://localhost:8080.

### Tokens Personalizaveis

Assim que a aplicação é iniciada, no pacote config são registados 5 tokens personalizaveis, sendo eles:

**TOKEN_1**
**TOKEN_2**
**TOKEN_3**
**TOKEN_4**
**TOKEN_5**

As configurações nas variáveis **LOCK_DURATION_SECONDS** e **BLOCK_DURATION_SECONDS** refletem para todos os tokens e IP's. A **LOCK_DURATION_SECONDS** significa o range de tempo que usaremos para controlar a quantidade de requisições e o **BLOCK_DURATION_SECONDS** é o tempo determinado que o IP ou Token ficará impossibilitado de realizar chamadas na API.


### Exemplos de Uso

curl `http://localhost:8080 -H "API_KEY: TOKEN_1"`.

### Fortio

docker exec -it <container_id> fortio load -c 2 -qps 12 -t 5s -H "API_KEY: TOKEN_2" http://app:8080

