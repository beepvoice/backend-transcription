# backend-transcription

Beep backend handling transcription of bites to text via Google Cloud.

## Quickstart

```
go build && ./backend-transcription
```

## Flags

Flags are supplied to the compiled go program in the form ```-flag=stuff```.

| Flag | Description | Default |
| ---- | ----------- | ------- |
| listen | Port number to listen on | 8080 |
| nats | URL of NATS | nats://localhost:4222 |
| api-key | Google Cloud API key | Something that works. Probably |

## API

## Scan Bites

```
GET /conversation/:key/scan
```

Get a list of transcription start times within a conversation key and specified timespan.

#### URL Params

| Name | Type | Description |
| ---- | ---- | ----------- |
| key | String | Audio transcription's bite's conversation's ID. |

#### Querystring

| Name | Type | Description |
| ---- | ---- | ----------- |
| from | Epoch timestamp | Time to start scanning from |
| to | Epoch timestamp | Time to scan to |

#### Success (200 OK)

```
Content-Type: application/json
```

```
{
  "previous": <Timestamp of transcription before <starts>>,
  "starts": [Timestamp, Timestamp...],
  "next": <Timestamp of transcription after <starts>>,
}
```

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Malformed input (from/to not timestamp, key not alphanumeric). |
| 500 | NATs or protobuf serialisation encountered errors. |

---

### Get Bite

```
GET /conversation/:key/start/:start
```

Get a specific ```transcription```.

#### URL Params

| Name | Type | Description |
| ---- | ---- | ----------- |
| key | String | Audio transcription's conversation's ID. |
| start | Epoch timestamp | Time the audio transcription starts. |

#### Success (200 OK)

Plaintext transcription

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | start is not an uint/key is not an alphanumeric string/specified bite could not be found |
| 500 | NATs or protobuf serialisation encountered errors. |
