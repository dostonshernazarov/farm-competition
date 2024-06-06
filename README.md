# farmish-crm


<h2>Clone project</h2>
<a href="https://github.com/dostonshernazarov/farm-competition">Github</a>

<h2>Create .env file. </h2>
<h3>Example: <a href="./.env">.example.env</a> </h3>

<h2>How to run</h2>

```
make migrate-up
```

**1. With make file** <br>

```
make run
```

**2. With terminal** <br>

```
go run cmd/app/main.go
```

**3. With docker compose** <br>

```
docker compose up
```

<h2><a href="https://www.postgresql.org/docs/current/datatype-json.html">*JSONB</a> type in project</h2>
```
[{"capacity":1, "time":14:00}, {"capacity":2, "time":15:00}, {"capacity":3, "time":16:00}]
```
