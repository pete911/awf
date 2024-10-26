# awf

Search aws resources. Search is done in local storage, you need to run `import` first.

## import AWS data

- `AWS_PROFILE=<profile> awf import`
- for all profiles `for p in $(aws configure list-profiles);do echo $p; AWS_PROFILE=$p awf import; done`
