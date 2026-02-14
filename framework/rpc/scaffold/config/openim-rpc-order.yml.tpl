rpc:
  registerName: openim-rpc-order
  rpc:
    listenIP: 0.0.0.0
    ports: [10300]
    registerIP: ""
    autoSetPorts: false
  prometheus:
    enable: true
    ports: [12300]
  rateLimiter:
    enable: true
    bucket: 2000
    window: 1000
    cpuThreshold: 900
  circuitBreaker:
    enable: true
    success: 0.6
    request: 20
    bucket: 10
    window: 3000
