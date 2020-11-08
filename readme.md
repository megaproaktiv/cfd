# `cfd` CloudFormation "one click" deploy

Generate a CloudFormation template with `goformation` and embedded Stack Management.

- Define your template in "template.go".
- Build program for your os
- call program `cfd` with
    - `cfd deploy` to deploy your stack to CloudFormation with the current credentials
    - `cfd destroy` to destroy your stack
    - `cfd status` to show the status of your stack
    - `cfd show` to show your stack




