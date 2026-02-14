api:
  api:
    listenIP: 0.0.0.0
    ports: [10002]
  prometheus:
    enable: true
    ports: [12002]
  rateLimiter:
    enable: true
    bucket: 1000
    window: 1000
    cpuThreshold: 900
