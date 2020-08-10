# Testing Prerequisites

## Required Objects

To execute the test suite, you will need the following objects set up in OneFuse:

- a Naming Sequence
- a Naming Policy that:
	- uses the above Naming Sequence in its Naming Template
- a Microsoft Module Endpoint

## Configuration Parameters

You will need to set the following values as appropriate for your environment:

- the connection information (scheme, address, port, username, and password) for your 
  OneFuse installation
- the database id of the Naming Policy above
- the databse id and name of the Microsoft Module Endpoint above

These configuration parameters are set as environment variables. The included 
"config.env" file contains the necessary exports. Simply provide values for each 
option appropriate to your environment and run "source config.env" in your terminal 
before running the tests. 

