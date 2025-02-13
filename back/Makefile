.PHONY:  create-keypair migrate-create migrate-up migrate-down migrate-force

PWD = $(shell pwd)
MPATH = $(PWD)/migrations
PORT = 5432


# Default number of migrations to execute up or down

N = 1

create-keypair:
    @echo "Creating an rsa 256 key pair"
    openssl genpkey -algorithm RSA -out $(PWD)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
    openssl rsa -in $(PWD)/rsa_private_$(ENV).pem -pubout -out $(PWD)/rsa_public_$(ENV).pem

migrate-create:
    @echo "---Creating migration files---"
    migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

migrate-up:
    -migrate -source file://$(MPATH) -database postgresql://admin:root@postgres_gulmarket:$(PORT)/postgres?sslmode=disable up || echo "No new migrations or an error occurred"

migrate-down:
    migrate -source file://$(MPATH) -database postgresql://admin:root@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
    migrate -source file://$(MPATH) -database postgresql://admin:root@postgres_gulmarket:$(PORT)/postgres?sslmode=disable force $(VERSION)



run:
    -go run ./
