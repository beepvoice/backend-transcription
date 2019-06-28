# backend-transcription

Beep backend handling transcription of bites to text via Google Cloud. Is completely within the backend and has no exposed endpoints.

## Quickstart

```
go build && ./backend-transcription
```

## Environment Variables

Supply environment variables by either exporting them or editing ```.env```.

| ENV | Description | Default |
| ---- | ----------- | ------- |
| NATS | Host and port of nats | nats://localhost:4222 |
| API_KEY | Google Cloud API key | Something that works. Probably |

