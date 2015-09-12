# hume

Hume is a service for performing empirical data validation tests against previously seen and "known good" values/distributions/etc. 

## API

URL | Operations | Action
--- | --- | ---
`/hume/api/config/metrics` | `GET` | List all metric types 
`/hume/api/config/metrics/{metric-type}` | `GET` | List all metrics of a given type and provide a template for a new one
`/hume/api/config/metrics/{metric-type}` | `POST` | Add a metrics 
`/hume/api/config/metrics/{metric-type}/{id}` | `GET`,`UPDATE`,`DELETE` | Get, Update, Delete a metric
`/hume/api/config/sources` | `GET` | List all source types
`/hume/api/config/sources/{source-type}` | `GET` | List all sources of a given type and provide a template for a new one
`/hume/api/config/sources/{source-type}` | `POST` | Add a source
`/hume/api/config/sources/{source-type}/{id}` | `GET`,`UPDATE`,`DELETE` | Get, Update, Delete a source
`/hume/api/job` | `GET` | List all recent jobs (paginated)
`/hume/api/job` | `POST` | Start a new job
`/hume/api/job/{id}` | `GET` | Get a detailed job status
`/hume/api/job/{id}` | `POST` | Post a new job status (used by agent)
`/hume/api/job/{id}` | `DELETE` | Stop a job

## Example Configuration
### Metrics
Metrics are configurable functions that can be applied to files of fields to output statsitics that can then be persisted and/or evaluated.

Metric names are globally-unique, so they can be easily referenced in a data source configuration, and are associated with a spcific metric type. Some metric types may require configuration while others do not. For example, a metric to calculate the number of fields of a data set (important for file processing):
```json
{
	"name":"Field Count",
	"type":"FieldCount"
}
```

But a metric for, say, checking the date format of a field would look like:
```json
{
	"name":"Date Format: YYYY-MM-DD",
	"type":"DateFormat",
	"format":"%Y-%m-%d"
}
```

### Sources
Sources define an input source, such as a file or a sql query, their expected fields, any metrics and associated evaluators.
```json
{
	"name":"Example File",
	"type":"File",
	"delimiter":"|",
	"fields":[
		"id",
		"value"
	],
	"has_header":True,
	"tests":[
		{
			"metric":"Field Count",
			"persist":True,
			"evaulate":[
				{
					"evaluator":"SingleValueMatch",
					"value":2
					"threshhold":1.00
				}
			]
		}
	]
}
```
