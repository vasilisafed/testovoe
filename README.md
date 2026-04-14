# Signature Service
- Go 1.20+

## Run

```bash
go run ./app
```

The server starts on `:8080`.

## API

### Create device

`POST /devices`

```json
{
  "id": "device-1",
  "algorithm": "RSA",
  "label": "Primary signer"
}
```

### List devices

`GET /devices`

### Get device

`GET /devices/{id}`

### Sign transaction

`POST /devices/{id}/sign`

```json
{
  "data": "transaction payload"
}
```

Response:

```json
{
  "signature": "<base64-signature>",
  "signed_data": "<counter>\n<data>\n<previous-signature-or-base64-device-id>"
}
```

## Design Notes

- Devices are stored in memory behind a repository interface so the storage can be replaced later.
- The signing logic is algorithm-agnostic through a signer factory, which makes adding new algorithms straightforward.
- Signature creation is serialized with a mutex to keep `signature_counter` strictly monotonically increasing in this single-node in-memory implementation.
- On the first signature, the previous signature value is `base64(device.id)` as required by the challenge.

## Verification

```bash
go test ./...
go build ./...
```

The automated tests cover:

- device creation and retrieval
- duplicate device rejection
- signature formatting and counter updates
- concurrent signing behavior
- main HTTP lifecycle endpoints
