# Negativações

Este serviço funciona como Façade para a aplicação legada de negativações, ajudando com uma parte do processamento nas requisições de buscas por CPF.

## Pré-requisitos

É necessário a instalação do docker, que pode ser feita seguindo o passo-a-passo do site oficial: https://docs.docker.com/get-docker/

Para a simulação completa do sistema (legado, banco de dados, api-gateway, serviço de negativações), é preciso instalar o docker-compose: https://docs.docker.com/compose/install/

## Subindo o serviço

Executando o comando `make init` serão carregadas todas as dependências de código, passando pela etapa de build, criação de imagem docker e por fim será executada a stack.

Ao final da inicialização você poderá acessar a API.

#### Sincronizando com a apliacação legada
```shell
curl -X POST -H "Api-Key: c5b6e72c-5b04-4bd2-ba5e-c85a253191dc" http://localhost/negativation/synchronize

HTTP/1.1 200 OK
```
```json
{"data":null,"status":"success"}
```

#### Buscando negativações por CPF
```shell
curl -X GET -H "Api-Key: c5b6e72c-5b04-4bd2-ba5e-c85a253191dc" http://localhost/negativation?cpf=515.374.764-67

HTTP/1.1 200 OK
```
```json
{
  "data":[{
    "companyDocument":"59291534000167",
    "companyName":"ABC S.A.",
    "customerDocument":"51537476467",
    "value":1235.23,
    "contract":"bc063153-fb9e-4334-9a6c-0d069a42065b",
    "debtDate":"2015-11-13T23:32:51Z",
    "inclusionDate":"2020-11-13T23:32:51Z"
  },{
    "companyDocument":"77723018000146",
    "companyName":"123 S.A.",
    "customerDocument":"51537476467",
    "value":400,
    "contract":"5f206825-3cfe-412f-8302-cc1b24a179b0",
    "debtDate":"2015-10-12T23:32:51Z",
    "inclusionDate":"2020-10-12T23:32:51Z"
  }],
  "status":"success"
}
```

## Executando testes

Para executar os testes é preciso instalar uma dependência para geração criptografia reversível: https://docs.cossacklabs.com/themis/installation/installation-from-packages/

Após instalado, basta executar `make test`