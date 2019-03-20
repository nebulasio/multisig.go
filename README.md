# Vote

## Build
```
make build
./nebms
```

**Command list.**

### create key & signature
```
create             Create a private key.
sign <file_path>   Sign the data file.
```

### contract
```
contract <contract_config_file_path>                 Create contract.
```

### data
```
data delete <address>                                Create 'delete manager' data.
data add <address>                                   Create 'add manager' data.
data replace <oldAddress> <newAddress>               Create 'replace manager' data.
data send_nas <txs_file_path>                        Create 'send nas' data.
data update_send_nas_rule <send_nas_rule_file_path>  Create 'update send nas signature rules' data.
data update_sys_config <sys_config_file_path>        Create 'update sys config' data.
data merge <data_file1_path> <data_file2_path> ...   Merge data file.
```
