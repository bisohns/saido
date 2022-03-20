(echo # config^
hosts:^
  children:^
    "127.0.0.1":^
      connection:^
        type: local^
metrics:^
- memory^
- cpu)>config-test.yaml
