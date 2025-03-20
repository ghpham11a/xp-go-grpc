# Usage

POST

```
{
  "id": "1",
  "email": "john.doe@example.com",
  "dateOfBirth": "1985-07-13",
  "accountNumber": "ACCT-12345-XYZ",
  "balance": 1234.56,
  "createdAt": "2025-03-07T10:25:00Z"
}
```

# 1. Initialize project

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

```
go mod init xp-go-grpc
go mod tidy
```

```
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```
go get google.golang.org/grpc
```

```
make sure Go sees where you installed the binaries
```

Compile protobuf files


Generate protobuf files for Client and server seperatly 

```
protoc \
  -I=. \
  -I=/Users/user/Documents/googleapis \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/accounts.proto
```

# 2. Run project

```
go run server/main.go
```

In another terminal

```
go run client/main.go --name="Luke"
```

###### 3. Setup PostgreSQL

Setup PostgreSQL containers using values from dev-postgres-values.yaml. This is in the k8s folder.

```sh
helm install xp-postgres -f dev-postgres-config.yaml bitnami/postgresql
```

Connect to the PostgreSQL in the container via psql. First set the password.

```sh
# set admin password with Powershell
$POSTGRES_PASSWORD = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String((kubectl get secret --namespace default xp-postgres-postgresql -o jsonpath="{.data.postgres-password}")))

# set admin password in bash
secret=$(kubectl get secret --namespace default xp-postgres-postgresql -o jsonpath="{.data.postgres-password}")

POSTGRES_PASSWORD=$(echo "$secret" | base64 --decode)
```

Using the password from above, get into the container and launch psql

```sh
kubectl run xp-postgres-postgresql-client --rm --tty -i --restart='Never' --namespace default --image docker.io/bitnami/postgresql:17.0.0-debian-12-r6 --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host xp-postgres-postgresql -U postgres -d postgres -p 5432
```

Create the table


```sh
postgres=# CREATE TABLE accounts (id SERIAL PRIMARY KEY, email VARCHAR(50) NOT NULL, date_of_birth DATE, account_number VARCHAR(20) UNIQUE, balance DECIMAL(18,2) DEFAULT 0.00, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
```

To check that it was created successfully

```sh
postgres=# \dt
```

Grant priveledges to user setup in dev-postgres-config.yaml

```
GRANT ALL PRIVILEGES ON TABLE accounts TO "postgres-username";
GRANT USAGE, SELECT, UPDATE ON SEQUENCE accounts_id_seq TO "postgres-username";
```

Quit psql

```sh
postgres=# \q
```