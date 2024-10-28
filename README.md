# awf

Search aws resources - AWFind. Search is done in local storage, you need to run `import` first. The difference between
aws cli and this tool is:
- speed, search is done on already imported resources
- no throttling (same as above, search is done against local storage)
- if the import is done for multiple accounts, search is done across all imported resources

## build/install

### download

-  download [release](https://github.com/pete911/awf/releases)

### brew

- add tap `brew tap pete911/tap`
- install `brew install awf`

## import AWS data

- `AWS_PROFILE=<profile> awf import`
- for all profiles `for p in $(aws configure list-profiles);do echo $p; AWS_PROFILE=$p awf import; done`

Imported resources are stored under `$HOME/.awf/` directory. In case import fails, or data needs to be cleaned up,
simply run `rm -r ~/.awf/*` and re-run the import.

## commands

Output columns are 'squashed' to 25 characters. If you see in the middle of the output `..`, it means it has been
'squashed'. If you need to see full length columns, use `--trim=false` flag. E.g. `aws subnet --trim=false 10.0.0.0/16`.

- network interfaces `aws ni <IP|CIDR|ID>` e.g. `aws ni 10.0.0.0/16` or `aws ni 10.60.3.25 10.5.0.0/24`
- network vpcs `aws vpc <IP|CIDR|ID>`
- network subnets `aws subnet <IP|CIDR|ID>`

## examples

```
aws ni 10.60.3.25 10.5.0.0/24

ACCOUNT ID    AWS PROFILE   ENI         TYPE  DESCRIPTION                PRIVATE IP    PUBLIC IP      VPC ID      VPC NAME  SUBNET ID    SUBNET NAME
123456789012  test-one      eni-xyzabc  alb   ELB app/test..ne/abcxyz12  10.60.3.25                   vpc-xyz123  test      subnet-deef  test-a
987654321098  test-two      eni-abcxyz  nat   Interface fo..NAT Gateway  10.5.0.1      216.58.212.238 vpc-123xyz  test      subnet-beef  test-a
987654321098  test-two      eni-abc123  nlb   ELB net/test..wo-alb-inte  10.5.0.2                     vpc-123456  test      subnet-geef  test-b
```

## release

Releases are published when a new tag is created e.g. `git tag -m "initial release" v0.0.1 && git push --follow-tags`.
