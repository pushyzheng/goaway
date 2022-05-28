## goaway

ga~~te~~way -> g**o**away

## Introduction

Goaway is a **lightweight** proxy server for protecting your website.

features:
1. Don't have any store engine, configure by yaml file.
2. Supports web server or file server
3. Auth and permissions
4. Prometheus metrics

## Quick start

Create config.yaml file in root directory firstly:

```yaml
server:
  port: 3000
  domain: example.com
  cookie-expired-hours: 72
  debug: true
  prometheus-path: /gateway/metrics
```