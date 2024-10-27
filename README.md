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

## examples

```
aws ni 10.60.3.25 10.5.0.0/24

ACCOUNT ID    AWS PROFILE   ENI         TYPE                       DESCRIPTION                               PRIVATE IP    PUBLIC IP      VPC         SUBNET
123456789012  test-one      eni-xyzabc  application_load_balancer  ELB app/test-one/abcxyz123                10.60.3.25                   vpc-xyz123  subnet-deef
987654321098  test-two      eni-abcxyz  nat_gateway                Interface for NAT Gateway nat-xyzabc1...  10.5.0.1      216.58.212.238 vpc-123xyz  subnet-beef
987654321098  test-two      eni-abc123  network_load_balancer      ELB net/test-two-alb-internal/xyzabc1...  10.5.0.2                     vpc-123456  subnet-geef
```

