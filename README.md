# awf

Search aws resources - AWFind. Search is done in local storage, you need to run `import` first. The difference between
aws cli and this tool is:
- speed, search is done on already imported resources
- no throttling (same as above, search is done against local storage)
- if the import is done for multiple accounts, search is done across all imported resources

## import AWS data

- `AWS_PROFILE=<profile> awf import`
- for all profiles `for p in $(aws configure list-profiles);do echo $p; AWS_PROFILE=$p awf import; done`

Imported resources are stored under `$HOME/.awf/` directory. In case import fails, or data needs to be cleaned up,
simply run `rm -r ~/.awf/*` and re-run the import.

## commands

- network interfaces `aws ni <IP|CIDR|ID>` e.g. `aws ni 10.0.0.0/16` or `aws ni 10.60.3.25 10.5.0.0/24`
- network vpcs `aws vpc <IP|CIDR|ID>`
- network subnets `aws subnet <IP|CIDR|ID>`

## examples

```
aws ni 10.60.3.25 10.5.0.0/24

ACCOUNT ID    AWS PROFILE   ENI         TYPE  DESCRIPTION                               PRIVATE IP    PUBLIC IP      VPC ID      VPC NAME  SUBNET ID    SUBNET NAME
123456789012  test-one      eni-xyzabc  alb   ELB app/test-one/abcxyz123                10.60.3.25                   vpc-xyz123  test      subnet-deef  test-a
987654321098  test-two      eni-abcxyz  nat   Interface for NAT Gateway nat-xyzabc1...  10.5.0.1      216.58.212.238 vpc-123xyz  test      subnet-beef  test-a
987654321098  test-two      eni-abc123  nlb   ELB net/test-two-alb-internal/xyzabc1...  10.5.0.2                     vpc-123456  test      subnet-geef  test-b
```
