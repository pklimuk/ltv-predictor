![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/pklimuk/ltv-predictor/ci.yml)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/pklimuk/ltv-predictor/codeql.yml?label=codeQL)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/pklimuk/ltv-predictor)
[![Coverage Status](https://coveralls.io/repos/github/pklimuk/ltv-predictor/badge.svg?branch=main)](https://coveralls.io/github/pklimuk/ltv-predictor?branch=main)

LTV predictor
---
### Description
This is a simple LTV predictor. It works with two types of input files(csv and json). Examples of input files structure could be found in the `testData` folder.

### Usage
To run the predictor you need to run the following command:
```
go run main.go -source <pathToSourceFile> [-model <model> -aggregate <aggregateByField> -predictionLength <predictionLength>]
```
Where:
```
source - path to the source file
```
```
model - predictor model. Could be one of the following: 
  -linearExtrapolation(default)
  -linearRegression
```
```
aggregate - field by which the data will be aggregated. Could be one of the following: 
  -country(default)
  -campaign
```
```
predictionLength - length of the prediction in days, default is 60
```