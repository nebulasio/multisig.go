# Nebulas Multi-Sign & Voting

## Build
```
make build
```

## Command-Line Toolset

### Generate Muli-Sign Smart Contract
Generate a new smart contract base on existing template.js -o contract-0310.js
```
neb-cli contract gen template.js contract.conf 
```
Generate a new smart contract with the config. contract-0310.js

### Deploy the smart contract
Generate a new smart contract base on existing template.js
```
neb-cli contract deploy contract-0301.js keystore.txt
```
Return the address of the deployed smart contract, the keystore.txt. can be any account.

### Create an account (cold environment)
Create an account with private key instead of keystore
```
neb-cli account create
```
Generate key.txt file， the format of a key.txt file:
<private_key>,<account_addr>

### Generate data file been signed
- Generate data file through a trans_file.csv
```
neb-cli data-gen <trans_file.csv> -o <file_content>
```
- Add a new signee address
```
neb-cli data-gen add-signee <address> -o <file_content>
```
- Remove a signee address
```
neb-cli data-gen remove-signee <address> -o <file_content>
```
- Replace a signee address
```
neb-cli data-gen replace-signee <address_origin> <address_new> -o <file_content>
```
- Update rules when sending transactions
```
neb-cli data-gen update-rules <rule_file> -o <file_content>
```

- Update the base configration for the smart contract
```
neb-cli data-gen update-constitution <config_file> -o <file_content>
```
- Merge multiple data need to be signed
```
neb-cli data-gen merge-file <directory>  -o <file_content>
```

#### Format for the files
- trans_file.csv
   ```
   NAS,n1Hb8rKQodFjQdksf8BDhbnnqgVrEDHLRC8,100.00
   NAS,n1dn3ZjJPRnL4UdePFkwQ5b9AJqPxnsQvM2,101.00
   NAS,n1ckLKYqUTSjgoEeV4HDtbTUg8pb38vmP68,88.88
   ```

### Sign content with private key (cold environment)
```
neb-cli sign <file_content> <key_file> -o signed-0311.txt
```
Generate a signed file: signed-0311.txt
