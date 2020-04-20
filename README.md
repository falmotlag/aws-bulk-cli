# aws bulk-cli
A thin wrapper around AWS CLI that enables running AWS CLI commands on multiple AWS accounts.

## Prerequisites
- [AWS CLI](https://aws.amazon.com/cli/)

## Configuration
### CLI Options


| Option                          | Description                                                |
| ----------------------------- | ---------------------------------------------------------- |
| ```accounts```                | Comma-delimited list of AWS accounts                       |
| ```cross-account-role-name``` | Name of cross-account IAM role used to access AWS accounts |
| ```config```                  | Path to ```aws bulk-cli``` config file.                          |
| ```region```                  | AWS region to run ```aws bulk-cli```                                                         |
|       ```profile```                        |    Named profile to use to run ```aws bulk-cli```                                                        |

### Config File
```aws bulk-cli```  can be configured with a json file. 
#### Config File Template

```json
{
    "aws": {
        "profile": "profile_name",
        "region": "region_name",
        "cross-account-role-name": "role_name",
        "accounts": [
            "000000000000", 
            "111111111111"
        ]
    }
}
```
## How it works
1. Iterates over the list of AWS accounts
2. Assume role ```arn:aws:iam::<account_id>:role/<cross-account-role-name>```
3. Run command in AWS account
4. Wait for each command in all accounts to return
5. Display command result 
