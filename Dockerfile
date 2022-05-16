####################################################################################################
# Step 1: Build the app
####################################################################################################

FROM rwynn/monstache-builder-cache-rel6:1.0.7 AS build-app

RUN mkdir /app

WORKDIR /app

RUN go mod download

COPY . .

RUN make release

RUN go build -buildmode=plugin -o build/plugin.so mapper_plugin.go

####################################################################################################
# Step 2: Copy output build file to an alpine image
####################################################################################################

FROM rwynn/monstache-alpine:3.15.0

ENTRYPOINT ["/bin/monstache"]

COPY --from=build-app /app/build/linux-amd64/monstache /bin/monstache

COPY --from=build-app /app/build/plugin.so /bin/plugin.so
